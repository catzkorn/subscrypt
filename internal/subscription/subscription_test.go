package subscription

import (
	"github.com/Catzkorn/subscrypt/internal/plaid"
	"github.com/shopspring/decimal"
	"reflect"
	"testing"
	"time"
)

func TestProcessTransactions(t *testing.T) {
	t.Run("Returns a list of subscriptions after processing a known subscription from the statement of transactions", func(t *testing.T) {
		transactions := plaid.TransactionList{Transactions: []plaid.Transaction{{Amount: 9.99, Date: "2020-09-12", Name: "Netflix"}}}
		amount, _ := decimal.NewFromString("9.99")
		want := []Subscription{{ID: 0, Name: "Netflix", Amount: amount, DateDue: time.Date(2020, time.Now().Month() + 1, 12, 0, 0, 0, 0, time.UTC)}}
		got := ProcessTransactions(transactions)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}

	})

	t.Run("Returns only a known subscription from the statement of transactions", func(t *testing.T) {
		transactions := plaid.TransactionList{Transactions: []plaid.Transaction{{Amount: 9.99, Date: "2020-09-12", Name: "Netflix"}, {Amount: 9.99, Date: "2020-09-12", Name: "Spotify"}}}
		amount, _ := decimal.NewFromString("9.99")
		want := []Subscription{{ID: 0, Name: "Netflix", Amount: amount, DateDue: time.Date(2020, time.Now().Month() + 1, 12, 0, 0, 0, 0, time.UTC)}}
		got := ProcessTransactions(transactions)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}

	})

	t.Run("Does not allow duplicate transactions", func(t *testing.T) {
		transactions := plaid.TransactionList{Transactions: []plaid.Transaction{{Amount: 9.99, Date: "2020-09-12", Name: "Netflix"}, {Amount: 9.99, Date: "2020-09-12", Name: "Spotify"}, {Amount: 9.99, Date: "2020-08-12", Name: "Netflix"}}}
		amount, _ := decimal.NewFromString("9.99")
		want := []Subscription{{ID: 0, Name: "Netflix", Amount: amount, DateDue: time.Date(2020, time.Now().Month() + 1, 12, 0, 0, 0, 0, time.UTC)}}
		got := ProcessTransactions(transactions)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}

	})
}