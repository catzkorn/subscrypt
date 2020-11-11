package main

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/shopspring/decimal"
)

func TestDatabaseConnection(t *testing.T) {

	t.Run("tests a successful database connection", func(t *testing.T) {
		_, err := NewSubscriptionStore("user=charlotte  host=localhost port=5432 database=subscryptdb sslmode=disable")

		if err != nil {
			t.Errorf("failed to create subscription store: %v", err)
		}
	})

	t.Run("tests a database connection failure", func(t *testing.T) {
		_, err := NewSubscriptionStore("gary the gopher")

		if err == nil {
			t.Errorf("connected to database that doesn't exist")
		}

	})
}

func TestDatabaseExecContext(t *testing.T) {
	store, err := NewSubscriptionStore("user=charlotte  host=localhost port=5432 database=subscryptdb sslmode=disable")

	if err != nil {
		t.Errorf("unexpected connection error: %v", err)
	}

	t.Run("adds a subscription and retrieves all added subscriptions", func(t *testing.T) {
		amount, _ := decimal.NewFromString("8.00")

		subscription := Subscription{
			Name:    "Netflix",
			Amount:  amount,
			DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC),
		}

		err := store.RecordSubscription(subscription)

		if err != nil {
			t.Errorf("unexpected database error: %v", err)
		}

		subscriptions, err := store.GetSubscriptions()

		if err != nil {
			t.Errorf("unexpected database error: %v", err)
		}

		if subscriptions[0].ID == 0 {
			t.Errorf("Database did not return an ID, got %v want %v", 0, subscriptions[0].ID)
		}

		if subscriptions[0].Name != subscription.Name {
			t.Errorf("Database did not return correct subscription name, got %s want %s", subscriptions[0].Name, subscription.Name)
		}

		if !subscriptions[0].Amount.Equal(subscription.Amount) {
			t.Errorf("Database did not return correct amount, got %#v want %#v", subscriptions[0].Amount, subscription.Amount)
		}

		if !subscriptions[0].DateDue.Equal(subscription.DateDue) {
			t.Errorf("Database did not return correct subscription date, got %s want %s", subscriptions[0].DateDue, subscription.DateDue)
		}

		clearSubscriptionsTable()

	})

}

func clearSubscriptionsTable() error {
	db, err := sql.Open("pgx", "user=charlotte  host=localhost port=5432 database=subscryptdb sslmode=disable")
	if err != nil {
		return fmt.Errorf("unexpected connection error: %w", err)
	}

	db.ExecContext(context.Background(), "TRUNCATE TABLE subscriptions;")

	return err

}
