package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Catzkorn/subscrypt/internal/database"
	"github.com/Catzkorn/subscrypt/internal/subscription"
	"github.com/shopspring/decimal"
)

func TestCreatingSubsAndRetrievingThem(t *testing.T) {
	store := database.NewInMemorySubscriptionStore()
	server := subscription.NewSubscriptionServer(store)
	amount, _ := decimal.NewFromString("100")
	subscription := subscription.Subscription{1, "Netflix", amount, time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC)}

	server.ServeHTTP(httptest.NewRecorder(), subscription.newPostSubscriptionRequest(subscription))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, subscription.NewGetSubscriptionRequest())
	subscription.assertStatus(t, response.Code, http.StatusOK)

	got := subscription.getSubscriptionsFromResponse(t, response.Body)
	subscription.assertSubscriptions(t, got, []subscription.Subscription{{1, "Netflix", amount, time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC)}})
}

func TestCreatingSubsAndRetrievingThemFromDatabase(t *testing.T) {
	store, _ := database.NewDatabaseConnection(database.DatabaseConnTestString)
	server := subscription.NewSubscriptionServer(store)
	amount, _ := decimal.NewFromString("100")
	subscription := subscription.Subscription{
		Name:    "Netflix",
		Amount:  amount,
		DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC),
	}

	server.ServeHTTP(httptest.NewRecorder(), subscription.newPostSubscriptionRequest(subscription))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, subscription.NewGetSubscriptionRequest())
	assertStatus(t, response.Code, http.StatusOK)

	got := subscription.getSubscriptionsFromResponse(t, response.Body)
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

	database.ClearSubscriptionsTable()
}
