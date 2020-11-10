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

		request := newSubscriptionRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getSubscriptionsFromResponse(t, response.Body)
		assertSubscriptions(t, got, wantedSubscriptions)

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

func newSubscriptionRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	return req
}