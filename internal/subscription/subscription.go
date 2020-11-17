package subscription

import (
	"time"

	"github.com/shopspring/decimal"
)

// Subscription defines a subscription. ID is unique per subscription.
// Name is the name of the subscription stored as a string.
// Amount is the cost of the subscription, stored as a decimal.
// DateDue is the date that the subscription is due on, stored as a date.
type Subscription struct {
	ID      int					`json:"id"`
	Name    string				`json:"name"`
	Amount  decimal.Decimal 	`json:"amount"`
	DateDue time.Time 			`json:"dateDue"`
}
