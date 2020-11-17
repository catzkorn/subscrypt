package integration_test

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Catzkorn/subscrypt/internal/database"
	"github.com/Catzkorn/subscrypt/internal/server"
	"github.com/Catzkorn/subscrypt/internal/subscription"
	"github.com/shopspring/decimal"
)

const indexTemplatePath = "../../web/index.html"

func TestCreatingSubsAndRetrievingThem(t *testing.T) {
	store := database.NewInMemorySubscriptionStore()
	testServer := server.NewServer(store, indexTemplatePath)
	amount, _ := decimal.NewFromString("100")
	wantedSubscriptions := []subscription.Subscription{
		{ID: 1, Name: "Netflix", Amount: amount, DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC)},
	}

	request := newPostFormRequest(url.Values{"name": {"Netflix"}, "amount": {"9.98"}, "date": {"2020-11-12"}})
	response := httptest.NewRecorder()
	testServer.ServeHTTP(response, request)

	request = newGetSubscriptionRequest()
	response = httptest.NewRecorder()

	testServer.ServeHTTP(response, request)
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
	}

	bodyString := string(body)
	got := bodyString

	res := strings.Contains(got, wantedSubscriptions[0].Name)

	if res != true {
		t.Errorf("webpage did not contain subscription of name %v", wantedSubscriptions[0].Name)
	}
}

func TestDeletingSubscriptionFromInMemoryStore(t *testing.T) {
	store := database.NewInMemorySubscriptionStore()
	testServer := server.NewServer(store, indexTemplatePath)

	amount, _ := decimal.NewFromString("100")
	subscription := subscription.Subscription{
		Name:    "Netflix",
		Amount:  amount,
		DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC),
	}
	storedSubscription, err := store.RecordSubscription(subscription)
	if err != nil {
		fmt.Println(err)
	}

	request := newDeleteSubscriptionRequest(storedSubscription.ID)
	response := httptest.NewRecorder()

	testServer.ServeHTTP(response, request)

	assertStatus(t, response.Code, http.StatusOK)

	gotSubscription, err := store.GetSubscription(storedSubscription.ID)
	if err != nil {
		fmt.Println(err)
	}

	if gotSubscription != nil {
		t.Errorf("subscription not deleted, got %v for given id, want nil", gotSubscription)
	}

	err = clearSubscriptionsTable()
	assertDatabaseError(t, err)
}

func TestCreatingSubsAndRetrievingThemFromDatabase(t *testing.T) {
	store, _ := database.NewDatabaseConnection(os.Getenv("DATABASE_CONN_STRING"))
	testServer := server.NewServer(store, indexTemplatePath)
	amount, _ := decimal.NewFromString("100")
	wantedSubscriptions := []subscription.Subscription{
		{ID: 1, Name: "Netflix", Amount: amount, DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC)},
	}

	request := newPostFormRequest(url.Values{"name": {"Netflix"}, "amount": {"9.98"}, "date": {"2020-11-12"}})
	response := httptest.NewRecorder()
	testServer.ServeHTTP(response, request)

	request = newGetSubscriptionRequest()
	response = httptest.NewRecorder()

	testServer.ServeHTTP(response, request)
	body, err := ioutil.ReadAll(response.Body)

	bodyString := string(body)
	got := bodyString

	fmt.Println(bodyString)

	res := strings.Contains(got, wantedSubscriptions[0].Name)

	if res != true {
		t.Errorf("webpage did not contain subscription of name %v", wantedSubscriptions[0].Name)
	}

	err = clearSubscriptionsTable()
	assertDatabaseError(t, err)
}

func TestDeletingSubscriptionFromDatabase(t *testing.T) {
	store, _ := database.NewDatabaseConnection(os.Getenv("DATABASE_CONN_STRING"))
	testServer := server.NewServer(store, indexTemplatePath)

	amount, _ := decimal.NewFromString("100")
	subscription := subscription.Subscription{
		Name:    "Netflix",
		Amount:  amount,
		DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC),
	}
	storedSubscription, err := store.RecordSubscription(subscription)
	if err != nil {
		fmt.Println(err)
	}

	request := newDeleteSubscriptionRequest(storedSubscription.ID)
	response := httptest.NewRecorder()

	testServer.ServeHTTP(response, request)

	assertStatus(t, response.Code, http.StatusOK)

	gotSubscription, err := store.GetSubscription(storedSubscription.ID)
	if err != nil {
		fmt.Println(err)
	}

	if gotSubscription != nil {
		t.Errorf("subscription not deleted, got %v for given id, want nil", gotSubscription)
	}

	err = clearSubscriptionsTable()
	assertDatabaseError(t, err)
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

func assertDatabaseError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected database error: %v", err)
	}
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func newGetSubscriptionRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	return req
}

func newPostFormRequest(url url.Values) *http.Request {
	var bodyStr = []byte(url.Encode())
	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyStr))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func newDeleteSubscriptionRequest(ID int) *http.Request {
	bodyStr := []byte(fmt.Sprintf("{\"id\": %v}", ID))
	url := fmt.Sprintf("/api/subscriptions/%v", ID)
	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(bodyStr))
	if err != nil {
		panic(err)
	}
	return req
}
