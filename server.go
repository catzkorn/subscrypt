package main

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"log"
	"net/http"
)

const JsonContentType = "application/json"

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
		s.processGetSubscription(w)
	case http.MethodPost:
		s.processPostSubscription(w, r)
	}
}

// processGetSubscription processes the GET subscription request, returning the store subscriptions as json
func (s *SubscriptionServer) processGetSubscription(w http.ResponseWriter) {
	w.Header().Set("content-type", JsonContentType)
	err := json.NewEncoder(w).Encode(s.store.GetSubscriptions())
	if err != nil {
		log.Fatalln(err)
	}
}

// processPostSubscription tells the SubscriptionStore to record the subscription from the post body
func (s *SubscriptionServer) processPostSubscription(w http.ResponseWriter, r *http.Request) {
	var subscription Subscription

	err := json.NewDecoder(r.Body).Decode(&subscription)
	if err != nil {
		log.Fatalln(err)
	}
	s.store.RecordSubscription(subscription)
	w.WriteHeader(http.StatusAccepted)
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
	ID int
	Name string
	Amount decimal.Decimal
	DateDue string
}