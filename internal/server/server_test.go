package server

import (
	"bytes"
	"fmt"
	"github.com/Catzkorn/subscrypt/internal/plaid"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/Catzkorn/subscrypt/internal/subscription"
	"github.com/shopspring/decimal"
)

type StubDataStore struct {
	subscriptions []subscription.Subscription
	deleteCount []int
}

func (s *StubDataStore) GetSubscriptions() ([]subscription.Subscription, error) {
	amount, _ := decimal.NewFromString("100.99")
	return []subscription.Subscription{{ID: 1, Name: "Netflix", Amount: amount, DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC)}}, nil
}

func (s *StubDataStore) RecordSubscription(subscription subscription.Subscription) (*subscription.Subscription, error) {
	s.subscriptions = append(s.subscriptions, subscription)
	return &subscription, nil
}

func (s *StubDataStore) GetSubscription(ID int) (*subscription.Subscription, error) {
	amount, _ := decimal.NewFromString("100.99")
	retrievedSubscription := subscription.Subscription{ID: 1, Name: "Netflix", Amount: amount, DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC)}
	if ID != 1 {
		return nil, nil
	}
	return &retrievedSubscription, nil
}

func (s *StubDataStore) DeleteSubscription(ID int) error {
	s.deleteCount = append(s.deleteCount, ID)
	return nil
}

type stubTransactionAPI struct{

}

func (s * stubTransactionAPI) GetTransactions() (plaid.TransactionList, error){
	transactions := plaid.TransactionList{Transactions: []plaid.Transaction{{Amount: 9.99, Date: "2020-09-12", MerchantName: "Netflix", Name: "Netflix"}}}
	return transactions, nil
}

func TestGetTransactions(t *testing.T) {
	t.Run("Successfully calls the transactionAPI", func(t *testing.T) {
		store := &StubDataStore{}
		transactionAPI := &stubTransactionAPI{}
		server := NewServer(store, transactionAPI)

		request, _ := http.NewRequest(http.MethodGet, "/api/transactions/", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)


	})
}


func TestGETSubscriptions(t *testing.T) {

	t.Run("return a subscription", func(t *testing.T) {
		amount, _ := decimal.NewFromString("100.99")
		wantedSubscriptions := []subscription.Subscription{
			{ID: 1, Name: "Netflix", Amount: amount, DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC)},
		}

		store := &StubDataStore{subscriptions: wantedSubscriptions}
		transactionAPI := &stubTransactionAPI{}
		server := NewServer(store, transactionAPI)

		request := newGetSubscriptionRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
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
	})
}

func TestStoreSubscription(t *testing.T) {

	t.Run("stores a subscription we POST to the server", func(t *testing.T) {
		store := &StubDataStore{}
		transactionAPI := &stubTransactionAPI{}
		server := NewServer(store, transactionAPI)


		request := newPostFormRequest(url.Values{"name": {"Netflix"}, "amount": {"9.98"}, "date": {"2020-11-12"}})

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		if len(store.subscriptions) != 1 {
			t.Errorf("got %d calls to RecordSubscription want %d", len(store.subscriptions), 1)
		}
	})
}

func TestDeleteSubscriptionAPI(t *testing.T) {

	t.Run("deletes the specified subscription from the data store and returns 200", func(t *testing.T) {
		subscriptions := []subscription.Subscription{{ID: 1}}
		store := &StubDataStore{subscriptions: subscriptions}
		transactionAPI := &stubTransactionAPI{}
		server := NewServer(store, transactionAPI)

		request := newDeleteSubscriptionRequest(1)

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		if len(store.deleteCount) != 1 {
			t.Errorf("got %d calls to DeleteSubscription want %d", len(store.deleteCount), 1)
		}
	})

	t.Run("returns 404 if given subscription ID doesn't exist", func(t *testing.T) {
		subscriptions := []subscription.Subscription{{ID: 1}}
		store := &StubDataStore{subscriptions: subscriptions}
		transactionAPI := &stubTransactionAPI{}
		server := NewServer(store, transactionAPI)

		request := newDeleteSubscriptionRequest(2)

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
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

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}