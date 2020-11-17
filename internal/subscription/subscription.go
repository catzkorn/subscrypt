package subscription

import (
	"github.com/Catzkorn/subscrypt/internal/plaid"
	"time"

	"github.com/shopspring/decimal"
)

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

func ProcessTransactions(transactions plaid.TransactionList) []Subscription {

	knownSubscriptions := []string{"Netflix", "Touchstone Climbing", "SparkFun"}

	var subscriptions []Subscription

	for _, transaction := range transactions.Transactions {

		if stringInSlice(transaction.Name, knownSubscriptions) {
			amount := decimal.NewFromFloat32(transaction.Amount)
			subscriptionDate := processDate(transaction.Date)
			subscription := Subscription{Name: transaction.Name, Amount: amount, DateDue: subscriptionDate}
			subscriptions = append(subscriptions, subscription)
		}
	}

	return subscriptions
}

func processDate(date string) time.Time {
	layout := "2006-01-02"
	str := date
	t, _ := time.Parse(layout, str)

	var subscriptionDate time.Time

	if t.Day() <= time.Now().Day() {
		subscriptionDate = time.Date(time.Now().Year(), time.Now().Month() + 1, t.Day(), 0, 0, 0, 0, time.UTC)
	} else {
		// next month date
		subscriptionDate = time.Date(time.Now().Year(), time.Now().Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	}

	return subscriptionDate
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}