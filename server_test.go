package main

import (
	"encoding/json"
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
	return []Subscription{{1, "Netflix", 100, "30"},}
}

func (s *StubSubscriptionStore) RecordSubscription(subscription Subscription) {
	s.subscriptions = append(s.subscriptions, subscription)
}

func TestGETSubscriptions(t *testing.T) {
	t.Run("Returns 200 OK", func(t *testing.T) {
		store := &StubSubscriptionStore{}
		server := NewSubscriptionServer(store)
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})
	
	t.Run("return a JSON of subscription", func(t *testing.T) {
		wantedSubscriptions := []Subscription{
			{1, "Netflix", 100, "30"},
		}

		store := StubSubscriptionStore{wantedSubscriptions}
		server := NewSubscriptionServer(&store)

		request := newGetSubscriptionRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getSubscriptionsFromResponse(t, response.Body)
		assertSubscriptions(t, got, wantedSubscriptions)

	})
}

func TestStoreSubscription(t *testing.T) {
	t.Run("returns 202 Accepted", func(t *testing.T) {
		store := &StubSubscriptionStore{}
		server := NewSubscriptionServer(store)
		subscription := Subscription{1, "Netflix", 100, "30"}
		request := newPostSubscriptionRequest(subscription)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)
	})

	t.Run("stores a subscription we POST to the server", func(t *testing.T) {
		subscription := Subscription{1, "Netflix", 100, "30"}

		store := &StubSubscriptionStore{}
		server := NewSubscriptionServer(store)

		request := newPostSubscriptionRequest(subscription)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		if len(store.subscriptions) != 1 {
			t.Errorf("got %d calls to RecordSubscription want %d", len(store.subscriptions), 1)
		}

		if store.subscriptions[0] != subscription {
			t.Errorf("did not store correct winner got %q want %q", store.subscriptions[0], subscription)
		}

	})
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
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

func assertSubscriptions(t *testing.T, got, want []Subscription) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func newGetSubscriptionRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	return req
}

func newPostSubscriptionRequest(subscription Subscription) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	return req
}