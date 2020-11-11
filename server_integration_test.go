package main

import (
	"github.com/shopspring/decimal"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatingSubsAndRetrievingThem(t *testing.T) {
	store := NewInMemorySubscriptionStore()
	server := NewSubscriptionServer(store)
	amount, _ := decimal.NewFromString("100")
	subscription := Subscription{1, "Netflix", amount, "30"}

	server.ServeHTTP(httptest.NewRecorder(), newPostSubscriptionRequest(subscription))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetSubscriptionRequest())
	assertStatus(t, response.Code, http.StatusOK)

	got := getSubscriptionsFromResponse(t, response.Body)
	assertSubscriptions(t, got, []Subscription{{1, "Netflix", amount, "30"}})
}