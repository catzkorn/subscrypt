package integration_test

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"github.com/Catzkorn/subscrypt/internal/database"
	"github.com/Catzkorn/subscrypt/internal/server"
	"github.com/Catzkorn/subscrypt/internal/subscription"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)
//
//func TestCreatingSubsAndRetrievingThem(t *testing.T) {
//	store := database.NewInMemorySubscriptionStore()
//	testServer := server.NewServer(store)
//	amount, _ := decimal.NewFromString("100")
//	wantedSubscriptions := []subscription.Subscription{
//		{ID: 1, Name: "Netflix", Amount: amount, DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC)},
//	}
//
//	request := newPostFormRequest(url.Values{"name": {"Netflix"}, "amount": {"9.98"}, "date": {"2020-11-12"}})
//	response := httptest.NewRecorder()
//	testServer.ServeHTTP(response, request)
//
//	request = newGetSubscriptionRequest()
//	response = httptest.NewRecorder()
//
//	testServer.ServeHTTP(response, request)
//	body, err := ioutil.ReadAll(response.Body)
//
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	bodyString := string(body)
//	got := bodyString
//
//	res := strings.Contains(got, wantedSubscriptions[0].Name)
//
//	if res != true {
//		t.Errorf("webpage did not contain subscription of name %v", wantedSubscriptions[0].Name)
//	}
//
//}

func TestCreatingSubsAndRetrievingThemFromDatabase(t *testing.T) {
	store, _ := database.NewDatabaseConnection(os.Getenv("DATABASE_CONN_STRING"))
	testServer := server.NewServer(store)
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

	err = clearSubscriptionsTable()
	assertDatabaseError(t, err)
}

func TestDeletingSubscriptionFromDatabase(t *testing.T) {
	store, _ := database.NewDatabaseConnection(os.Getenv("DATABASE_CONN_STRING"))
	testServer := server.NewServer(store)

	amount, _ := decimal.NewFromString("100")
	subscription := subscription.Subscription{
		ID: 1,
		Name: "Netflix",
		Amount: amount,
		DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC),
	}
	storedSubscription, err := store.RecordSubscription(subscription)

	request := newPostDeleteRequest(url.Values{"ID": {fmt.Sprint(storedSubscription.ID)}})
	response := httptest.NewRecorder()

	testServer.ServeHTTP(response, request)
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
	}

	bodyString := string(body)
	got := bodyString

	res := !strings.Contains(got, storedSubscription.Name)

	if res != true {
		t.Errorf("subscription not deleted, webpage contained subscription of name %v", storedSubscription.Name)
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

func newPostDeleteRequest(url url.Values) *http.Request {
	var bodyStr = []byte(url.Encode())
	req, err := http.NewRequest(http.MethodPost, "/delete", bytes.NewBuffer(bodyStr))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req
}