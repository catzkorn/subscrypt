package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Catzkorn/subscrypt/internal/subscription"

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
func (d *Database) RecordSubscription(sub subscription.Subscription) (*subscription.Subscription, error) {
	var id int
	var name string
	var amount pgtype.Numeric
	var dateDue time.Time

	err := d.database.QueryRowContext(context.Background(), "INSERT INTO subscriptions (name, amount, date_due) VALUES ($1, $2, $3) RETURNING id, name, amount, date_due", sub.Name, sub.Amount, sub.DateDue).Scan(&id, &name, &amount, &dateDue)
	if err != nil {
		return nil, fmt.Errorf("unexpected insert error: %w", err)
	}

	switch {
	case err == sql.ErrNoRows:
		return nil, fmt.Errorf("no user with id %d\n", id)
	case err != nil:
		return nil, fmt.Errorf("query error: %v\n", err)
	default:
		newSubscription := subscription.Subscription{
			ID:      id,
			Name:    name,
			Amount:  decimal.NewFromBigInt(amount.Int, amount.Exp),
			DateDue: dateDue,
		}
		return &newSubscription, nil
	}

}

// GetSubscriptions retrieves all subscriptions from the subscription database
func (d *Database) GetSubscriptions() ([]subscription.Subscription, error) {
	rows, err := d.database.QueryContext(context.Background(), "SELECT * FROM subscriptions;")
	if err != nil {
		return nil, fmt.Errorf("unexpected retrieve error: %w", err)
	}

	var subscriptions []subscription.Subscription

	for rows.Next() {
		var id int
		var name string
		var amount pgtype.Numeric
		var dateDue time.Time
		if err := rows.Scan(&id, &name, &amount, &dateDue); err != nil {
			log.Fatal(err)
		}
		subscriptions = append(subscriptions, subscription.Subscription{
			ID:      id,
			Name:    name,
			Amount:  decimal.NewFromBigInt(amount.Int, amount.Exp),
			DateDue: dateDue,
		})
	}
	return subscriptions, nil
}

func (d *Database) DeleteSubscription(subscriptionID int) error {
	_, err := d.database.ExecContext(context.Background(), "DELETE FROM subscriptions WHERE id = $1;", subscriptionID)
	if err != nil {
		return fmt.Errorf("unexpected insert error: %w", err)
	}

	return nil
}
