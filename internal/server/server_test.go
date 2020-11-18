package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
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
	userprofile := userprofile.Userprofile{Name: name, Email: email}

	return &userprofile, nil
}
func (s *StubDataStore) GetUserDetails() (*userprofile.Userprofile, error) {
	userprofile := userprofile.Userprofile{}

	return &userprofile, nil
}

type stubTransactionAPI struct {
}

func (s * stubTransactionAPI) GetTransactions() (plaid.TransactionList, error){
	transactions := plaid.TransactionList{Transactions: []plaid.Transaction{{Amount: 9.99, Date: "2020-09-12", Name: "Netflix"}}}
	return transactions, nil
}

func TestGetTransactions(t *testing.T) {
	t.Run("Successfully calls the transactionAPI", func(t *testing.T) {
		store := &StubDataStore{}
		transactionAPI := &stubTransactionAPI{}
		server := NewServer(store, indexTemplatePath, &StubMailer{}, transactionAPI)

		request, _ := http.NewRequest(http.MethodGet, "/api/transactions/", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusFound)
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
		server := NewServer(store, indexTemplatePath, &StubMailer{}, transactionAPI)

		request := newGetSubscriptionRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getSubscriptionsFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assertSubscriptions(t, got, wantedSubscriptions)
		assertContentType(t, response, JSONContentType)
	})
}

func TestStoreSubscription(t *testing.T) {

	t.Run("stores a subscription we POST to the server", func(t *testing.T) {
		amount, _ := decimal.NewFromString("100.99")
		subscription := subscription.Subscription{Name: "Netflix", Amount: amount, DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC)}

		store := &StubDataStore{}

		transactionAPI := &stubTransactionAPI{}
		server := NewServer(store, indexTemplatePath, &StubMailer{}, transactionAPI)

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
		server := NewServer(store, indexTemplatePath, &StubMailer{}, transactionAPI)

		request := newPostReminderRequest(t, subscriptions[0].ID)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		fmt.Println(response.Body.String())
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
		t.Fatalf("Unable to parse response from server %q into slice of Subscription, '%v'", body, err)
	}

	return
}

func TestDeleteSubscriptionAPI(t *testing.T) {

	t.Run("deletes the specified subscription from the data store and returns 200", func(t *testing.T) {
		subscriptions := []subscription.Subscription{{ID: 1}}
		store := &StubDataStore{subscriptions: subscriptions}

		transactionAPI := &stubTransactionAPI{}
		server := NewServer(store, indexTemplatePath, &StubMailer{}, transactionAPI)

		request := newDeleteSubscriptionRequest(1)

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
		server := NewServer(store, indexTemplatePath, &StubMailer{}, transactionAPI)

		request := newDeleteSubscriptionRequest(2)

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestUserHandler(t *testing.T) {

	t.Run("tests creation of a user", func(t *testing.T) {
		userprofile := userprofile.Userprofile{}

		store := &StubDataStore{userprofile: userprofile}

		transactionAPI := &stubTransactionAPI{}
		server := NewServer(store, indexTemplatePath, &StubMailer{}, transactionAPI)

		request := newPostUserRequest(t, "Charlotte", os.Getenv("EMAIL"))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

	})

}

func newGetSubscriptionRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/api/subscriptions", nil)
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

func newDeleteSubscriptionRequest(ID int) *http.Request {
	bodyStr := []byte(fmt.Sprintf("{\"id\": %v}", ID))
	url := fmt.Sprintf("/api/subscriptions/%v", ID)
	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(bodyStr))
	if err != nil {
		panic(err)
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
