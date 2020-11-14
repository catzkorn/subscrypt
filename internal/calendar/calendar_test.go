package calendar

import (
	"os"
	"testing"
	"time"

	"github.com/Catzkorn/subscrypt/internal/reminder"
	"github.com/Catzkorn/subscrypt/internal/subscription"
	ics "github.com/arran4/golang-ical"
	"github.com/shopspring/decimal"
)

func TestCalendarInviteCreation(t *testing.T) {

	amount, _ := decimal.NewFromString("8.00")
	subscription := subscription.Subscription{
		ID:      1,
		Name:    "Netflix",
		Amount:  amount,
		DateDue: time.Date(2020, time.November, 16, 0, 0, 0, 0, time.UTC),
	}

	reminder := reminder.Reminder{
		Email:          os.Getenv("EMAIL"),
		SubscriptionID: 1,
		ReminderDate:   time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC),
	}

	t.Run("checks that a .ics file is containing the correct information", func(t *testing.T) {

		cal := CreateReminderInvite(subscription, reminder)

		t.Log(len(cal.Components))
		for _, component := range cal.Components {

			switch c := component.(type) {
			case *ics.VEvent:
				t.Log(c)
			default:
				t.Errorf("did not create a VEvent")
			}
		}

		// calfile := strings.NewReader(cal.Serialize())

		// _, err := ics.ParseCalendar(calfile)

	})

}
