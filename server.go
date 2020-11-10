package main

import (
	"encoding/json"
	"net/http"
)

// NewSubscriptionServer returns a instance of a SubscriptionServer
func NewSubscriptionServer(store SubscriptionStore) *SubscriptionServer {
	s := new(SubscriptionServer)
	s.store = store
	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(s.subscriptionHandler))
	s.Handler = router
	return s
}

// subscriptionHandler handles the routing logic for the index
func (s *SubscriptionServer) subscriptionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(s.store.GetSubscriptions())
	case http.MethodPost:
		var subscription Subscription

		err := json.NewDecoder(r.Body).Decode(&subscription)
		if err != nil {
			//TODO: log and return error
		}
		s.store.RecordSubscription(subscription)
		w.WriteHeader(http.StatusAccepted)
	}
}

// SubscriptionServer is the HTTP interface for subscription information
type SubscriptionServer struct {
	store SubscriptionStore
	http.Handler
}

// SubscriptionStore stores information about individual subscriptions
type SubscriptionStore interface {
	GetSubscriptions() []Subscription
	RecordSubscription(subscription Subscription)
}

// Subscription stores the id, name, amount and datedue of an individual subscription
type Subscription struct {
	Id int
	Name string
	Amount int
	DateDue string
}