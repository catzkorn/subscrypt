package server

import (
	"bytes"
	"fmt"
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

func (s *StubDataStore) DeleteSubscription(ID int) error {
	s.deleteCount = append(s.deleteCount, ID)
	return nil
}

func TestGETSubscriptions(t *testing.T) {

	t.Run("return a subscription", func(t *testing.T) {
		amount, _ := decimal.NewFromString("100.99")
		wantedSubscriptions := []subscription.Subscription{
			{ID: 1, Name: "Netflix", Amount: amount, DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC)},
		}

		store := &StubDataStore{subscriptions: wantedSubscriptions}
		server := NewServer(store)

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
		server := NewServer(store)


		request := newPostFormRequest(url.Values{"name": {"Netflix"}, "amount": {"9.98"}, "date": {"2020-11-12"}})

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		if len(store.subscriptions) != 1 {
			t.Errorf("got %d calls to RecordSubscription want %d", len(store.subscriptions), 1)
		}
	})
}

func TestDeleteSubscription(t *testing.T) {

	t.Run("deletes the specified subscription from the data store", func(t *testing.T) {
		subscriptions := []subscription.Subscription{{ID: 1}}
		store := &StubDataStore{subscriptions: subscriptions}
		server := NewServer(store)

		request := newPostDeleteRequest(url.Values{"ID": {"1"}})

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		if len(store.deleteCount) != 1 {
			t.Errorf("got %d calls to DeleteSubscription want %d", len(store.deleteCount), 1)
		}
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

func newPostDeleteRequest(url url.Values) *http.Request {
	var bodyStr = []byte(url.Encode())
	req, err := http.NewRequest(http.MethodPost, "/delete", bytes.NewBuffer(bodyStr))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req
}