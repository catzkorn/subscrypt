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
	"github.com/Catzkorn/subscrypt/internal/plaid"
	"github.com/Catzkorn/subscrypt/internal/server"
	"github.com/Catzkorn/subscrypt/internal/subscription"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/shopspring/decimal"
)

const indexTemplatePath = "../../web/index.html"
const JSONContentType = "application/json"

type StubMailer struct {
	sentEmail *mail.SGMailV3
}

func (s *StubMailer) Send(email *mail.SGMailV3) (*rest.Response, error) {
	s.sentEmail = email
	return &rest.Response{StatusCode: http.StatusAccepted}, nil
}

func TestCreatingSubsAndRetrievingThem(t *testing.T) {
	store := database.NewInMemorySubscriptionStore()

	api := &plaid.PlaidAPI{}

	testServer := server.NewServer(store, indexTemplatePath, &StubMailer{}, api)

	amount, _ := decimal.NewFromString("100")
	newSubscription := subscription.Subscription{
		Name:    "Netflix",
		Amount:  amount,
		DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC),
	}

	postRequest := newPostSubscriptionRequest(t, newSubscription)
	response := httptest.NewRecorder()
	testServer.ServeHTTP(response, postRequest)

	assertStatus(t, response.Code, http.StatusOK)

	getRequest := newGetSubscriptionRequest()
	response = httptest.NewRecorder()
	testServer.ServeHTTP(response, getRequest)

	got := getSubscriptionsFromResponse(t, response.Body)

	assertStatus(t, response.Code, http.StatusOK)
	assertSubscription(t, got[0], newSubscription)
	assertContentType(t, response, JSONContentType)
}

func TestDeletingSubscriptionFromInMemoryStore(t *testing.T) {
	store := database.NewInMemorySubscriptionStore()

	api := &plaid.PlaidAPI{}
	testServer := server.NewServer(store, indexTemplatePath, &StubMailer{}, api)

	amount, _ := decimal.NewFromString("100")
	newSubscription := subscription.Subscription{
		Name:    "Netflix",
		Amount:  amount,
		DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC),
	}
	storedSubscription, err := store.RecordSubscription(newSubscription)
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

	api := &plaid.PlaidAPI{}
	testServer := server.NewServer(store, indexTemplatePath, &StubMailer{}, api)

	amount, _ := decimal.NewFromString("100")
	newSubscription := subscription.Subscription{
		Name:    "Netflix",
		Amount:  amount,
		DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC),
	}

	postRequest := newPostSubscriptionRequest(t, newSubscription)
	response := httptest.NewRecorder()
	testServer.ServeHTTP(response, postRequest)

	assertStatus(t, response.Code, http.StatusOK)

	getRequest := newGetSubscriptionRequest()
	response = httptest.NewRecorder()
	testServer.ServeHTTP(response, getRequest)

	got := getSubscriptionsFromResponse(t, response.Body)

	assertStatus(t, response.Code, http.StatusOK)
	assertSubscription(t, got[0], newSubscription)
	assertContentType(t, response, JSONContentType)

	if got[0].Name != newSubscription.Name {
		t.Errorf("Subscription not saved and retrieved successfully, got ")
	}
}

func TestDeletingSubscriptionFromDatabase(t *testing.T) {
	store, _ := database.NewDatabaseConnection(os.Getenv("DATABASE_CONN_STRING"))

	api := &plaid.PlaidAPI{}
	testServer := server.NewServer(store, indexTemplatePath, &StubMailer{}, api)

	amount, _ := decimal.NewFromString("100")
	newSubscription := subscription.Subscription{
		Name:    "Netflix",
		Amount:  amount,
		DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC),
	}
	storedSubscription, err := store.RecordSubscription(newSubscription)
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

// Assertion Test Helpers

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

// assertSubscription checks that the Name, Amount and DateDue of the got and want subscriptions match
// It doesn't check the ID value
func assertSubscription(t *testing.T, got, want subscription.Subscription) {
	t.Helper()

	if got.Name != want.Name {
		t.Errorf("subscriptions not correct - Name mismatch, got %v want %v", got, want)
	} else if !reflect.DeepEqual(got.Amount, want.Amount) {
		t.Errorf("subscriptions not correct - Amount mismatch, got %v want %v", got, want)
	} else if got.DateDue != want.DateDue {
		t.Errorf("subscriptions not correct - DateDue mismatch, got %v want %v", got, want)
	}
}

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}

// New Request Test Methods

func newGetSubscriptionRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/api/subscriptions", nil)
	return req
}

func newPostSubscriptionRequest(t *testing.T, subscription subscription.Subscription) *http.Request {
	postBody, _ := json.Marshal(subscription)
	req, err := http.NewRequest(http.MethodPost, "/api/subscriptions", bytes.NewBuffer(postBody))

	if err != nil {
		t.Errorf("failed to generate new POST subscription request")
	}

	return req
}

func newDeleteSubscriptionRequest(ID int) *http.Request {
	bodyStr := []byte(fmt.Sprintf("{\"id\": %v}", ID))
	deleteURL := fmt.Sprintf("/api/subscriptions/%v", ID)
	req, err := http.NewRequest(http.MethodDelete, deleteURL, bytes.NewBuffer(bodyStr))
	if err != nil {
		panic(err)
	}
	return req
}

func getSubscriptionsFromResponse(t *testing.T, body io.Reader) (subscriptions []subscription.Subscription) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&subscriptions)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Subscription, '%v'", body, err)
	}

	return
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