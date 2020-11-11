package main

import (
	"bytes"
	"encoding/json"
	"github.com/shopspring/decimal"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type StubSubscriptionStore struct {
	subscriptions []Subscription
}

func (s *StubSubscriptionStore) GetSubscriptions() []Subscription {
	amount, _ := decimal.NewFromString("100.99")
	return []Subscription{{1, "Netflix", amount, "30"},}
}

func (s *StubSubscriptionStore) RecordSubscription(subscription Subscription) {
	s.subscriptions = append(s.subscriptions, subscription)
}

func TestGETSubscriptions(t *testing.T) {

	t.Run("return a JSON of subscription", func(t *testing.T) {
		amount, _ := decimal.NewFromString("100.99")
		wantedSubscriptions := []Subscription{
			{1, "Netflix", amount, "30"},
		}

		store := StubSubscriptionStore{wantedSubscriptions}
		server := NewSubscriptionServer(&store)

		request := newGetSubscriptionRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getSubscriptionsFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assertSubscriptions(t, got, wantedSubscriptions)
		assertContentType(t, response, JsonContentType)
	})
}

func TestStoreSubscription(t *testing.T) {

	t.Run("stores a subscription we POST to the server", func(t *testing.T) {
		amount, _ := decimal.NewFromString("100.99")
		subscription := Subscription{1, "Netflix", amount, "30"}

		store := &StubSubscriptionStore{}
		server := NewSubscriptionServer(store)

		request := newPostSubscriptionRequest(subscription)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.subscriptions) != 1 {
			t.Errorf("got %d calls to RecordSubscription want %d", len(store.subscriptions), 1)
		}

		if !reflect.DeepEqual(store.subscriptions[0], subscription) {
			t.Errorf("did not store correct winner got %v want %v", store.subscriptions[0], subscription)
		}
	})
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertSubscriptions(t *testing.T, got, want []Subscription) {
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

func getSubscriptionsFromResponse(t *testing.T, body io.Reader) (subscriptions []Subscription) {
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

func newPostSubscriptionRequest(subscription Subscription) *http.Request {
	postBody, _ := json.Marshal(subscription)
	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(postBody))

	return req
}