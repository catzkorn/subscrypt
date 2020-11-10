package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var got []Subscription

		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("Unable to parse response from server %q into slice of Subscription, '%v'", response.Body, err)
		}

	})
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}