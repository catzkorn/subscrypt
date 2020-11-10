package main

func NewInMemorySubscriptionStore() *InMemorySubscriptionStore {
	return &InMemorySubscriptionStore{[]Subscription{}}
}

type InMemorySubscriptionStore struct{
	subscriptions []Subscription
}

func (i *InMemorySubscriptionStore) GetSubscriptions() []Subscription {
	return i.subscriptions
}

func (i *InMemorySubscriptionStore) RecordSubscription(subscription Subscription) {
	i.subscriptions = append(i.subscriptions, subscription)
}