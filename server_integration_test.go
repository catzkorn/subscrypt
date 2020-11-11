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

func TestCreatingSubsAndRetrievingThemFromDatabase(t *testing.T) {
	store, _ := NewDatabaseConnection(DatabaseConnTestString)
	server := NewSubscriptionServer(store)
	amount, _ := decimal.NewFromString("100")
	subscription := Subscription{
		Name:    "Netflix",
		Amount:  amount,
		DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC),
	}

	server.ServeHTTP(httptest.NewRecorder(), newPostSubscriptionRequest(subscription))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetSubscriptionRequest())
	assertStatus(t, response.Code, http.StatusOK)

	got := getSubscriptionsFromResponse(t, response.Body)
	if got[0].ID == 0 {
		t.Errorf("Database did not return an ID, got %v want %v", 0, got[0].ID)
	}

	if got[0].Name != subscription.Name {
		t.Errorf("Database did not return correct subscription name, got %s want %s", got[0].Name, subscription.Name)
	}

	if !got[0].Amount.Equal(subscription.Amount) {
		t.Errorf("Database did not return correct amount, got %#v want %#v", got[0].Amount, subscription.Amount)
	}

	if !got[0].DateDue.Equal(subscription.DateDue) {
		t.Errorf("Database did not return correct subscription date, got %s want %s", got[0].DateDue, subscription.DateDue)
	}

	clearSubscriptionsTable()
}
