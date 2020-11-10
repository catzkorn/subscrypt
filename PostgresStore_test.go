package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/jackc/pgtype"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/shopspring/decimal"
)

type SubscriptionMock struct {
	ID      int
	Name    string
	Amount  decimal.Decimal
	DateDue time.Time
}

func TestDatabaseConnection(t *testing.T) {

	_, err := NewPostgresStore("user=charlotte  host=localhost port=5432 database=subscryptdb sslmode=disable")

	if err != nil {
		t.Errorf("failed to create postgres store: %v", err)
	}

}

func TestDatabaseExecContext(t *testing.T) {
	db, err := sql.Open("pgx", "user=charlotte  host=localhost port=5432 database=subscryptdb sslmode=disable")

	if err != nil {
		t.Errorf("unexpected connection error: %v", err)
	}

	t.Run("add subscription", func(t *testing.T) {

		subscription := SubscriptionMock{
			Name:    "Netflix",
			Amount:  decimal.NewFromFloat(7.99),
			DateDue: time.Now(),
		}

		_, err := db.ExecContext(context.Background(), "INSERT INTO subscriptions (name, amount, date_due) VALUES ($1, $2, $3)", subscription.Name, subscription.Amount, subscription.DateDue)

		if err != nil {
			t.Errorf("unexpected database error: %v", err)
		}
	})

	t.Run("get subscriptions", func(t *testing.T) {

		subscription := SubscriptionMock{
			Name:    "Netflix",
			Amount:  decimal.NewFromFloat(7.99),
			DateDue: time.Now(),
		}

		_, err := db.ExecContext(context.Background(), "INSERT INTO subscriptions (name, amount, date_due) VALUES ($1, $2, $3)", subscription.Name, subscription.Amount, subscription.DateDue)

		if err != nil {
			t.Errorf("unexpected database error: %v", err)
		}

		rows, err := db.QueryContext(context.Background(), "SELECT * FROM subscriptions;")

		if err != nil {
			t.Errorf("unexpected database error: %v", err)
		}

		var subscriptions []SubscriptionMock

		for rows.Next() {
			var id int
			var name string
			var amount pgtype.Numeric
			var dateDue time.Time
			if err := rows.Scan(&id, &name, &amount, &dateDue); err != nil {
				log.Fatal(err)
			}
			subscriptions = append(subscriptions, SubscriptionMock{
				ID:      id,
				Name:    name,
				Amount:  decimal.NewFromBigInt(amount.Int, amount.Exp),
				DateDue: dateDue,
			})
		}

		fmt.Printf("%v", subscriptions)
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
