package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Catzkorn/subscrypt/internal/subscription"
	"github.com/Catzkorn/subscrypt/internal/userprofile"
	"github.com/shopspring/decimal"
)

const indexTemplatePath = "../../web/index.html"

type StubDataStore struct {
	subscriptions []subscription.Subscription
	deleteCount   []int
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

func TestGETSubscriptions(t *testing.T) {

	t.Run("return a subscription", func(t *testing.T) {
		amount, _ := decimal.NewFromString("100.99")
		wantedSubscriptions := []subscription.Subscription{
			{ID: 1, Name: "Netflix", Amount: amount, DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC)},
		}

		store := &StubDataStore{subscriptions: wantedSubscriptions}
		server := NewServer(store, indexTemplatePath)

		request := newGetSubscriptionRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		body, err := ioutil.ReadAll(response.Body)

		if err != nil {
			fmt.Println(err)
		}

		bodyString := string(body)
		got := bodyString

		res := strings.Contains(got, wantedSubscriptions[0].Name)

		if res != true {
			t.Errorf("webpage did not contain subscription of name %v", wantedSubscriptions[0].Name)
		}
	})
}

func TestStoreSubscription(t *testing.T) {

	t.Run("stores a subscription we POST to the server", func(t *testing.T) {
		store := &StubDataStore{}
		server := NewServer(store, indexTemplatePath)

		request := newPostFormRequest(url.Values{"name": {"Netflix"}, "amount": {"9.98"}, "date": {"2020-11-12"}})

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		if len(store.subscriptions) != 1 {
			t.Errorf("got %d calls to RecordSubscription want %d", len(store.subscriptions), 1)
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
		server := NewServer(store, indexTemplatePath)

		request := newPostReminderRequest("test@test.com", fmt.Sprint(subscriptions[0].ID), "2020-11-13")
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
		t.Fatalf("Unable to parse response from server %q into slice of Subscription, '%v'", body, err)
	}

	return

}
func TestDeleteSubscriptionAPI(t *testing.T) {

	t.Run("deletes the specified subscription from the data store and returns 200", func(t *testing.T) {
		subscriptions := []subscription.Subscription{{ID: 1}}
		store := &StubDataStore{subscriptions: subscriptions}
		server := NewServer(store, indexTemplatePath)

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
		server := NewServer(store, indexTemplatePath)

		request := newDeleteSubscriptionRequest(2)

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func newGetSubscriptionRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	return req
}

func newPostFormRequest(url url.Values) *http.Request {
	var bodyStr = []byte(url.Encode())
	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyStr))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func newPostReminderRequest(email string, subscriptionID string, reminderDate string) *http.Request {
	urlValues := url.Values{"email": {email}, "subscriptionID": {subscriptionID}, "reminderDate": {reminderDate}}
	var bodyStr = []byte(urlValues.Encode())
	req, err := http.NewRequest(http.MethodPost, "/reminder", bytes.NewBuffer(bodyStr))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		panic(err)
	}
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
