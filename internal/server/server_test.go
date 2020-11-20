package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/Catzkorn/subscrypt/internal/plaid"

	"github.com/Catzkorn/subscrypt/internal/subscription"
	"github.com/Catzkorn/subscrypt/internal/userprofile"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/shopspring/decimal"
)

const indexTemplatePath = "../../web/index.html"

type StubMailer struct {
	sentEmail *mail.SGMailV3
}

func (s *StubMailer) Send(email *mail.SGMailV3) (*rest.Response, error) {
	s.sentEmail = email
	return &rest.Response{StatusCode: http.StatusAccepted}, nil
}

type StubDataStore struct {
	subscriptions []subscription.Subscription
	deleteCount   []int
	userprofile   userprofile.Userprofile
}

func (s *StubDataStore) GetSubscriptions() ([]subscription.Subscription, error) {
	amount, _ := decimal.NewFromString("100.99")
	return []subscription.Subscription{{ID: 1, Name: "Netflix", Amount: amount, DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC)}}, nil
}

func (s *StubDataStore) RecordSubscription(subscription subscription.Subscription) (*subscription.Subscription, error) {
	s.subscriptions = append(s.subscriptions, subscription)
	return &subscription, nil
}

func (s *StubDataStore) GetSubscription(ID int) (*subscription.Subscription, error) {
	amount, _ := decimal.NewFromString("100.99")
	retrievedSubscription := subscription.Subscription{ID: 1, Name: "Netflix", Amount: amount, DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC)}
	if ID != 1 {
		return nil, nil
	}
	return &retrievedSubscription, nil
}

func (s *StubDataStore) DeleteSubscription(ID int) error {
	s.deleteCount = append(s.deleteCount, ID)
	return nil
}

func (s *StubDataStore) RecordUserDetails(name string, email string) (*userprofile.Userprofile, error) {
	s.userprofile = userprofile.Userprofile{Name: name, Email: email}

	return &s.userprofile, nil
}
func (s *StubDataStore) GetUserDetails() (*userprofile.Userprofile, error) {

	return &s.userprofile, nil
}

type stubTransactionAPI struct {
	transactionCount   int
}

func (s *stubTransactionAPI) GetTransactions() (plaid.TransactionList, error) {
	transactions := plaid.TransactionList{Transactions: []plaid.Transaction{{Amount: 9.99, Date: "2020-09-12", Name: "Netflix"}}}
	s.transactionCount ++
	return transactions, nil
}

func TestGetTransactions(t *testing.T) {
	t.Run("Successfully calls the transactionAPI", func(t *testing.T) {
		store := &StubDataStore{}
		transactionAPI := &stubTransactionAPI{}
		server := NewServer(store, &StubMailer{}, transactionAPI)

		request, _ := http.NewRequest(http.MethodGet, "/api/transactions", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("Successfully loads transactions page", func(t *testing.T) {
		store := &StubDataStore{}
		transactionAPI := &stubTransactionAPI{}
		server := NewServer(store, &StubMailer{}, transactionAPI)

		request, _ := http.NewRequest(http.MethodGet, "/transactions/", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		fmt.Println(response.Body)

		assertStatus(t, response.Code, http.StatusOK)
	})
}

func TestLoadSubscriptions(t *testing.T) {
	t.Run("Successfully calls the transactionAPI and loads Transactions", func(t *testing.T) {
		store := &StubDataStore{}
		transactionAPI := &stubTransactionAPI{}
		server := NewServer(store, &StubMailer{}, transactionAPI)

		request, _ := http.NewRequest(http.MethodPost, "/api/transactions/load-subscriptions", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		if transactionAPI.transactionCount != 1 {
			t.Errorf("got %d calls to GetTransactions want %d", transactionAPI.transactionCount, 1)
		}

		assertStatus(t, response.Code, http.StatusOK)
	})
}

func TestGETSubscriptions(t *testing.T) {

	t.Run("return subscriptions in JSON format", func(t *testing.T) {
		amount, _ := decimal.NewFromString("100.99")
		wantedSubscriptions := []subscription.Subscription{
			{ID: 1, Name: "Netflix", Amount: amount, DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC)},
		}

		store := &StubDataStore{subscriptions: wantedSubscriptions}

		transactionAPI := &stubTransactionAPI{}
		server := NewServer(store, &StubMailer{}, transactionAPI)

		request := newGetSubscriptionRequest(t)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getSubscriptionsFromResponse(t, response.Body)
		assertSubscriptions(t, got, wantedSubscriptions)

		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, JSONContentType)
	})
}

func TestStoreSubscription(t *testing.T) {

	t.Run("stores a subscription we POST to the server", func(t *testing.T) {
		amount, _ := decimal.NewFromString("100.99")
		subscription := subscription.Subscription{Name: "Netflix", Amount: amount, DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC)}

		store := &StubDataStore{}
		transactionAPI := &stubTransactionAPI{}
		server := NewServer(store, &StubMailer{}, transactionAPI)

		request := newPostSubscriptionRequest(t, subscription)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusOK)

		if len(store.subscriptions) != 1 {
			t.Errorf("got %d calls to RecordSubscription want %d", len(store.subscriptions), 1)
		}

		if !reflect.DeepEqual(store.subscriptions[0], subscription) {
			t.Errorf("did not store correct subscription got %v want %v", store.subscriptions[0], subscription)
		}
	})
}

func TestCreateReminder(t *testing.T) {

	t.Run("creates a reminder for subscription and returns a confirmation that a reminder invite has been sent", func(t *testing.T) {
		amount, _ := decimal.NewFromString("100.99")
		subscriptions := []subscription.Subscription{
			{ID: 1, Name: "Netflix", Amount: amount, DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC)},
		}

		store := &StubDataStore{subscriptions: subscriptions}
		transactionAPI := &stubTransactionAPI{}
		server := NewServer(store, &StubMailer{}, transactionAPI)

		request := newPostReminderRequest(t, subscriptions[0].ID)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusOK)
	})

}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertSubscriptions(t *testing.T, got, want []subscription.Subscription) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}

func getSubscriptionsFromResponse(t *testing.T, body io.Reader) (subscriptions []subscription.Subscription) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&subscriptions)
	if err != nil {
		t.Fatalf("unable to parse response from server %q into slice of Subscription, '%v'", body, err)
	}
	return
}

func TestDeleteSubscriptionAPI(t *testing.T) {

	t.Run("deletes the specified subscription from the data store and returns 200", func(t *testing.T) {
		subscriptions := []subscription.Subscription{{ID: 1}}
		store := &StubDataStore{subscriptions: subscriptions}

		transactionAPI := &stubTransactionAPI{}
		server := NewServer(store, &StubMailer{}, transactionAPI)

		request := newDeleteSubscriptionRequest(t, 1)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		if len(store.deleteCount) != 1 {
			t.Errorf("got %d calls to DeleteSubscription want %d", len(store.deleteCount), 1)
		}
	})

	t.Run("returns 404 if given subscription ID doesn't exist", func(t *testing.T) {
		subscriptions := []subscription.Subscription{{ID: 1}}
		store := &StubDataStore{subscriptions: subscriptions}

		transactionAPI := &stubTransactionAPI{}
		server := NewServer(store, &StubMailer{}, transactionAPI)

		request := newDeleteSubscriptionRequest(t, 2)

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestUserHandler(t *testing.T) {

	t.Run("tests creation of a user", func(t *testing.T) {
		store := &StubDataStore{}
		transactionAPI := &stubTransactionAPI{}
		server := NewServer(store, &StubMailer{}, transactionAPI)

		request := newPostUserRequest(t, "Gary Gopher", "gary@gopher.com")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusOK)

		if store.userprofile.Name != "Gary Gopher" {
			t.Errorf("incorrect name set got %v want %v", store.userprofile.Name, "Gary Gopher")
		}

		if store.userprofile.Email != "gary@gopher.com" {
			t.Errorf("incorrect name set got %v want %v", store.userprofile.Email, "gary@gopher.com")
		}
	})

	t.Run("tests retrival of user", func(t *testing.T) {
		userProfile := userprofile.Userprofile{
			Name:  "Gary Gopher",
			Email: "gary@gopher.com",
		}

		store := &StubDataStore{userprofile: userProfile}
		transactionAPI := &stubTransactionAPI{}
		server := NewServer(store, &StubMailer{}, transactionAPI)

		request := newGetUserRequest(t)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusOK)

		var retrievedUserProfile userprofile.Userprofile

		err := json.NewDecoder(response.Body).Decode(&retrievedUserProfile)
		if err != nil {
			t.Fatalf("unable to parse response from server %q into user profile, '%v'", response.Body, err)
		}

		if retrievedUserProfile.Name != userProfile.Name {
			t.Errorf("incorrect name retrieved got %v want %v", retrievedUserProfile.Name, userProfile.Name)
		}

		if retrievedUserProfile.Email != userProfile.Email {
			t.Errorf("incorrect name retrieved got %v want %v", retrievedUserProfile.Email, userProfile.Email)
		}

		assertContentType(t, response, JSONContentType)
	})

}

func newGetSubscriptionRequest(t testing.TB) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, "/api/subscriptions", nil)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	return req
}

func newPostSubscriptionRequest(t *testing.T, subscription subscription.Subscription) *http.Request {
	postBody, _ := json.Marshal(subscription)
	req, err := http.NewRequest(http.MethodPost, "/api/subscriptions", bytes.NewBuffer(postBody))
	if err != nil {
		t.Errorf("failed to generate new POST subscription request")
	}
	return req
}

func newPostReminderRequest(t testing.TB, id int) *http.Request {
	t.Helper()
	subscription := subscription.Subscription{
		ID: id,
	}
	bodyStr, err := json.Marshal(&subscription)
	if err != nil {
		t.Fatalf("fail to marshal subscription: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/api/reminders", bytes.NewBuffer(bodyStr))
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	return req
}

func newDeleteSubscriptionRequest(t testing.TB, ID int) *http.Request {
	url := fmt.Sprintf("/api/subscriptions/%v", ID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	return req
}

func newPostUserRequest(t testing.TB, name string, email string) *http.Request {
	t.Helper()
	userProfile := userprofile.Userprofile{
		Name:  name,
		Email: email,
	}

	bodyStr, err := json.Marshal(&userProfile)
	if err != nil {
		t.Fatalf("fail to marshal user information: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(bodyStr))
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	return req
}

func newGetUserRequest(t testing.TB) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, "/api/users", nil)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	return req
}
