package database

import "github.com/Catzkorn/subscrypt/internal/subscription"

// NewInMemorySubscriptionStore returns a instance of InMemorySubscriptionStore
func NewInMemorySubscriptionStore() *InMemorySubscriptionStore {
	return &InMemorySubscriptionStore{[]subscription.Subscription{}}
}

// InMemorySubscriptionStore stores information about individual subscriptions
type InMemorySubscriptionStore struct {
	subscriptions []subscription.Subscription
}

// GetSubscriptions is a method that returns all subscriptions
func (i *InMemorySubscriptionStore) GetSubscriptions() ([]subscription.Subscription, error) {
	return i.subscriptions, nil
}

// RecordSubscription is a method that stores a subscription into the store

func (i *InMemorySubscriptionStore) RecordSubscription(subscription subscription.Subscription) (*subscription.Subscription, error) {
	i.subscriptions = append(i.subscriptions, subscription)
	return &subscription, nil
}
