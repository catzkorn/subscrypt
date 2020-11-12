package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Catzkorn/subscrypt/internal/subscription"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/shopspring/decimal"
)

func TestDatabaseConnection(t *testing.T) {

	t.Run("tests a successful database connection", func(t *testing.T) {
		_, err := NewDatabaseConnection(os.Getenv("DATABASE_CONN_STRING"))
		assertDatabaseError(t, err)
	})

	t.Run("tests a database connection failure", func(t *testing.T) {
		_, err := NewDatabaseConnection("gary the gopher")

		if err == nil {
			t.Errorf("connected to database that doesn't exist")
		}

	})
}

func TestAddingSubscriptionToDB(t *testing.T) {
	store, err := NewDatabaseConnection(os.Getenv("DATABASE_CONN_STRING"))
	assertDatabaseError(t, err)

	t.Run("adds a subscription and retrieves all added subscriptions", func(t *testing.T) {
		wantedSubscription := createTestSubscription()

		subscription, err := store.RecordSubscription(wantedSubscription)
		assertDatabaseError(t, err)

		if subscription.ID == 0 {
			t.Errorf("Database did not return an ID, got %v want %v", 0, subscription.ID)
		}

		if subscription.Name != wantedSubscription.Name {
			t.Errorf("Database did not return correct subscription name, got %s want %s", subscription.Name, wantedSubscription.Name)
		}

		if !subscription.Amount.Equal(subscription.Amount) {
			t.Errorf("Database did not return correct amount, got %#v want %#v", subscription.Amount, wantedSubscription.Amount)
		}

		if !subscription.DateDue.Equal(wantedSubscription.DateDue) {
			t.Errorf("Database did not return correct subscription date, got %s want %s", subscription.DateDue, wantedSubscription.DateDue)
		}

		clearSubscriptionsTable()
	})

}

func TestDeletingSubscriptionFromDB(t *testing.T) {
	store, err := NewDatabaseConnection(os.Getenv("DATABASE_CONN_STRING"))
	assertDatabaseError(t, err)

	t.Run("deletes the subscription from the database", func(t *testing.T) {
		subscription := createTestSubscription()

		gotSubscription, err := store.RecordSubscription(subscription)
		assertDatabaseError(t, err)

		subscriptionID := gotSubscription.ID

		err = store.DeleteSubscription(subscriptionID)
		assertDatabaseError(t, err)

		subscriptions, err := store.GetSubscriptions()
		assertDatabaseError(t, err)

		if len(subscriptions) != 0 {
			t.Errorf("database did not delete subscription, got %v, wanted no subscriptions", subscriptions)
		}

		clearSubscriptionsTable()
	})
}

func createTestSubscription() subscription.Subscription {
	amount, _ := decimal.NewFromString("8.00")
	subscription := subscription.Subscription{
		Name:    "Netflix",
		Amount:  amount,
		DateDue: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC),
	}
	return subscription
}

func clearSubscriptionsTable() error {
	db, err := sql.Open("pgx", os.Getenv("DATABASE_CONN_STRING"))
	if err != nil {
		return fmt.Errorf("unexpected connection error: %w", err)
	}

	db.ExecContext(context.Background(), "TRUNCATE TABLE subscriptions;")

	return err
}

func assertDatabaseError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("unexpected database error: %v", err)
	}
}
