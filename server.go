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
	json.NewEncoder(w).Encode(s.store.GetSubscriptions())
}

type SubscriptionServer struct {
	store SubscriptionStore
	http.Handler
}

type SubscriptionStore interface {
	GetSubscriptions() []Subscription
}

type Subscription struct {
	Id int
	Name string
	Amount int
	DateDue string
}