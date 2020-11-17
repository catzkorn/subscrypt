package subscription

import (
	"github.com/Catzkorn/subscrypt/internal/plaid"
	"github.com/shopspring/decimal"
	"reflect"
	"testing"
	"time"
)

func TestProcessTransactions(t *testing.T) {
	t.Run("Returns a list of subscriptions after processing statement transactions", func(t *testing.T) {
		transactions := plaid.TransactionList{Transactions: []plaid.Transaction{{Amount: 9.99, Date: "2020-09-12", Name: "Netflix"}, {Amount: 4.99, Date: "2020-09-12", Name: "Spotify"}}}
		amount, _ := decimal.NewFromString("9.99")
		want := []Subscription{{ID: 1, Name: "Netflix", Amount: amount, DateDue: time.Date(2020, time.October, 12, 0, 0, 0, 0, time.UTC)}}
		got := ProcessTransactions(transactions)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}

	})
}