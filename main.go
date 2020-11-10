package main

import (
	"log"
	"net/http"
)

type InMemorySubscriptionStore struct{
	subscriptions []Subscription
}

func (i *InMemorySubscriptionStore) GetSubscriptions() []Subscription {
	return []Subscription{{1, "Netflixy", 100, "30"},}
}

func (i *InMemorySubscriptionStore) RecordSubscription(subscription Subscription) {
	i.subscriptions = append(i.subscriptions, subscription)
}

func main() {
	server := NewSubscriptionServer(&InMemorySubscriptionStore{})
	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}