package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatingSubsAndRetrievingThem(t *testing.T) {
	store := NewInMemorySubscriptionStore()
	server := NewSubscriptionServer(store)
	subscription := Subscription{1, "Netflix", "100", "30"}

	server.ServeHTTP(httptest.NewRecorder(), newPostSubscriptionRequest(subscription))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetSubscriptionRequest())
	assertStatus(t, response.Code, http.StatusOK)

	got := getSubscriptionsFromResponse(t, response.Body)
	assertSubscriptions(t, got, []Subscription{{1, "Netflix", "100", "30"}})
}