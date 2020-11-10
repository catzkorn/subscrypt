package main

import (
	"encoding/json"
	"net/http"
)

func NewSubscriptionServer(store SubscriptionStore) *SubscriptionServer {
	s := new(SubscriptionServer)

	s.store = store

	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(s.subscriptionHandler))

	s.Handler = router

	return s
}

func (s *SubscriptionServer) subscriptionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(s.store.GetSubscriptions())
	case http.MethodPost:
		subscription := Subscription{1, "Netflix", 100, "30"}
		s.store.RecordSubscription(subscription)
		w.WriteHeader(http.StatusAccepted)
	}
}


type SubscriptionServer struct {
	store SubscriptionStore
	http.Handler
}

type SubscriptionStore interface {
	GetSubscriptions() []Subscription
	RecordSubscription(subscription Subscription)
}

type Subscription struct {
	Id int
	Name string
	Amount int
	DateDue string
}