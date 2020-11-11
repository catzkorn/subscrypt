package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestCreatingSubsAndRetrievingThem(t *testing.T) {
	store := NewInMemorySubscriptionStore()
	server := NewSubscriptionServer(store)
	amount, _ := decimal.NewFromString("100")
	subscription := Subscription{1, "Netflix", amount, time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC)}

	server.ServeHTTP(httptest.NewRecorder(), newPostSubscriptionRequest(subscription))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetSubscriptionRequest())
	assertStatus(t, response.Code, http.StatusOK)

	got := getSubscriptionsFromResponse(t, response.Body)
	assertSubscriptions(t, got, []Subscription{{1, "Netflix", amount, time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC)}})
}
