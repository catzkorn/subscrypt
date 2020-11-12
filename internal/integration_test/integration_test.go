package integration_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/Catzkorn/subscrypt/internal/database"
	"github.com/Catzkorn/subscrypt/internal/server"
	"github.com/Catzkorn/subscrypt/internal/subscription"
	"github.com/shopspring/decimal"
)

func TestCreatingSubsAndRetrievingThem(t *testing.T) {
	store := database.NewInMemorySubscriptionStore()
	server := server.NewServer(store)
	amount, _ := decimal.NewFromString("100")
	subscriptionFML := subscription.Subscription{ID: 1, Name: "Netflix", Amount: amount, DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC)}

	server.ServeHTTP(httptest.NewRecorder(), newPostSubscriptionRequest(subscriptionFML))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetSubscriptionRequest())
	assertStatus(t, response.Code, http.StatusOK)

	got := getSubscriptionsFromResponse(t, response.Body)
	assertSubscriptions(t, got, []subscription.Subscription{{1, "Netflix", amount, time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC)}})
}

func TestCreatingSubsAndRetrievingThemFromDatabase(t *testing.T) {
	store, _ := database.NewDatabaseConnection(os.Getenv("DATABASE_CONN_STRING"))
	server := server.NewServer(store)
	amount, _ := decimal.NewFromString("100")
	subscriptionFML := subscription.Subscription{
		Name:    "Netflix",
		Amount:  amount,
		DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC),
	}

	server.ServeHTTP(httptest.NewRecorder(), newPostSubscriptionRequest(subscriptionFML))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetSubscriptionRequest())
	assertStatus(t, response.Code, http.StatusOK)

	got := getSubscriptionsFromResponse(t, response.Body)
	if got[0].ID == 0 {
		t.Errorf("Database did not return an ID, got %v want %v", 0, got[0].ID)
	}

	if got[0].Name != subscriptionFML.Name {
		t.Errorf("Database did not return correct subscription name, got %s want %s", got[0].Name, subscriptionFML.Name)
	}

	if !got[0].Amount.Equal(subscriptionFML.Amount) {
		t.Errorf("Database did not return correct amount, got %#v want %#v", got[0].Amount, subscriptionFML.Amount)
	}

	if !got[0].DateDue.Equal(subscriptionFML.DateDue) {
		t.Errorf("Database did not return correct subscription date, got %s want %s", got[0].DateDue, subscriptionFML.DateDue)
	}

	clearSubscriptionsTable()
}

func clearSubscriptionsTable() error {
	db, err := sql.Open("pgx", os.Getenv("DATABASE_CONN_STRING"))
	if err != nil {
		return fmt.Errorf("unexpected connection error: %w", err)
	}
	_, err = db.ExecContext(context.Background(), "TRUNCATE TABLE subscriptions;")
	if err != nil {
		return fmt.Errorf("unexpected connection error: %w", err)
	}

	return err
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

func newPostSubscriptionRequest(subscription subscription.Subscription) *http.Request {
	postBody, _ := json.Marshal(subscription)
	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(postBody))

	return req
}
