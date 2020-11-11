package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/shopspring/decimal"
)

// SubscriptionStore allows the user to store and read back subscriptions
type SubscriptionStore struct {
	database *sql.DB
}

// Subscription defines a subscription. ID is unique per subscription.
// Name is the name of the subscription stored as a string.
// Amount is the cost of the subscription, stored as a decimal.
// DateDue is the date that the subscription is due on, stored as a date.
type Subscription struct {
	ID      int
	Name    string
	Amount  decimal.Decimal
	DateDue time.Time
}

// NewSubscriptionStore starts connection with database
func NewSubscriptionStore(databaseDSN string) (*SubscriptionStore, error) {
	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		return nil, fmt.Errorf("unexpected connection error: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("unexpected connection error: %w", err)
	}

	return &SubscriptionStore{database: db}, nil
}

// RecordSubscription inserts a subscription into the subscription database
func (s *SubscriptionStore) RecordSubscription(subscription Subscription) error {
	_, err := s.database.ExecContext(context.Background(), "INSERT INTO subscriptions (name, amount, date_due) VALUES ($1, $2, $3)", subscription.Name, subscription.Amount, subscription.DateDue)
	if err != nil {
		return fmt.Errorf("unexpected insert error: %w", err)
	}

	return nil

}

// GetSubscriptions retrieves all subscriptions from the subscription database
func (s *SubscriptionStore) GetSubscriptions() ([]Subscription, error) {
	rows, err := s.database.QueryContext(context.Background(), "SELECT * FROM subscriptions;")
	if err != nil {
		return nil, fmt.Errorf("unexpected retrieve error: %w", err)
	}

	var subscriptions []Subscription

	for rows.Next() {
		var id int
		var name string
		var amount pgtype.Numeric
		var dateDue time.Time
		if err := rows.Scan(&id, &name, &amount, &dateDue); err != nil {
			log.Fatal(err)
		}
		subscriptions = append(subscriptions, Subscription{
			ID:      id,
			Name:    name,
			Amount:  decimal.NewFromBigInt(amount.Int, amount.Exp),
			DateDue: dateDue,
		})
	}
	return subscriptions, nil
}
