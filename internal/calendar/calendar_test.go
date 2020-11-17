package calendar

import (
	"fmt"
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
		if len(cal.Components) != 1 {
			t.Errorf("did not have the expected number of components got %v, want %v", cal.Components, 1)
		}

		event, ok := cal.Components[0].(*ics.VEvent)

		if !ok {
			t.Errorf("did not create a VEvent")
		}

		var uid string
		var attendee string

		for _, property := range event.Properties {

			switch property.IANAToken {
			case string(ics.PropertyUid):
				uid = property.Value
			case string(ics.PropertyAttendee):
				attendee = property.Value
				t.Logf(attendee)
			}
		}

		if uid != fmt.Sprintf("%v@subscrypt.com", subscription.ID) {
			t.Errorf("incorrect UID got %v want %v@subscrypt.com", uid, subscription.ID)
		}

		if attendee != fmt.Sprintf("mailto:%v", reminder.Email) {
			t.Errorf("incorrect attendee got %v want mailto:%v", attendee, reminder.Email)
		}
	})

}
