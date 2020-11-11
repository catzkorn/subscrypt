package main

// NewInMemorySubscriptionStore returns a instance of InMemorySubscriptionStore
func NewInMemorySubscriptionStore() *InMemorySubscriptionStore {
	return &InMemorySubscriptionStore{[]Subscription{}}
}

// InMemorySubscriptionStore stores information about individual subscriptions
type InMemorySubscriptionStore struct {
	subscriptions []Subscription
}

// GetSubscriptions is a method that returns all subscriptions
func (i *InMemorySubscriptionStore) GetSubscriptions() ([]Subscription, error) {
	return i.subscriptions, nil
}

// RecordSubscription is a method that stores a subscription into the store
func (i *InMemorySubscriptionStore) RecordSubscription(subscription Subscription) error {
	i.subscriptions = append(i.subscriptions, subscription)
	return nil
}
