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

// Database allows the user to store and read back subscriptions
type Database struct {
	database *sql.DB
}

// NewDatabaseConnection starts connection with database
func NewDatabaseConnection(databaseDSN string) (*Database, error) {
	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		return nil, fmt.Errorf("unexpected connection error: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("unexpected connection error: %w", err)
	}

	return &Database{database: db}, nil
}

// RecordSubscription inserts a subscription into the subscription database
func (d *Database) RecordSubscription(subscription Subscription) error {
	_, err := d.database.ExecContext(context.Background(), "INSERT INTO subscriptions (name, amount, date_due) VALUES ($1, $2, $3)", subscription.Name, subscription.Amount, subscription.DateDue)
	if err != nil {
		return fmt.Errorf("unexpected insert error: %w", err)
	}

	return nil

}

// GetSubscriptions retrieves all subscriptions from the subscription database
func (d *Database) GetSubscriptions() ([]Subscription, error) {
	rows, err := d.database.QueryContext(context.Background(), "SELECT * FROM subscriptions;")
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
