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
	"github.com/shopspring/decimal"
)

type StubDataStore struct {
	subscriptions []subscription.Subscription
}

func (s *StubDataStore) GetSubscriptions() ([]subscription.Subscription, error) {
	amount, _ := decimal.NewFromString("100.99")
	return []subscription.Subscription{{ID: 1, Name: "Netflix", Amount: amount, DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC)}}, nil
}

func (s *StubDataStore) RecordSubscription(subscription subscription.Subscription) (*subscription.Subscription, error) {
	s.subscriptions = append(s.subscriptions, subscription)
	return &subscription, nil
}

func TestGETSubscriptions(t *testing.T) {

	t.Run("return a subscription", func(t *testing.T) {
		amount, _ := decimal.NewFromString("100.99")
		wantedSubscriptions := []subscription.Subscription{
			{ID: 1, Name: "Netflix", Amount: amount, DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC)},
		}

		store := StubDataStore{wantedSubscriptions}
		server := NewServer(&store)

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
		server := NewServer(store)

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

		store := &StubDataStore{subscriptions}
		server := NewServer(store)

		request := newPostReminderRequest("test@test.com", subscriptions[0].ID, time.Date(2020, time.November, 6, 0, 0, 0, 0, time.UTC))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusAccepted)

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

func newPostReminderRequest(email string, subscriptionID int, reminderDate time.Time) *http.Request {
	var bodyStr = []byte(fmt.Sprintf("email=%s&subscriptionID=%v&date=%v", email, subscriptionID, reminderDate))
	req, err := http.NewRequest(http.MethodPost, "/reminder", bytes.NewBuffer(bodyStr))
	if err != nil {
		panic(err)
	}
	return req
}
