package email

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Catzkorn/subscrypt/internal/calendar"
	"github.com/Catzkorn/subscrypt/internal/reminder"
	"github.com/Catzkorn/subscrypt/internal/subscription"
	"github.com/Catzkorn/subscrypt/internal/userprofile"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/shopspring/decimal"
)

type StubMailer struct {
	sentEmail *mail.SGMailV3
}

func (s *StubMailer) Send(email *mail.SGMailV3) (*rest.Response, error) {
	s.sentEmail = email
	return &rest.Response{StatusCode: http.StatusAccepted}, nil
}

type StubDataStore struct {
	subscription subscription.Subscription
}

func (s *StubDataStore) GetSubscription(subscriptionID int) (*subscription.Subscription, error) {
	return &s.subscription, nil
}

func TestSendingAnEmail(t *testing.T) {
	reminder := reminder.Reminder{
		Email:          os.Getenv("EMAIL"),
		SubscriptionID: 1,
		ReminderDate:   time.Date(2020, time.December, 11, 0, 0, 0, 0, time.UTC),
	}

	amount, _ := decimal.NewFromString("8.00")
	subscription := subscription.Subscription{
		ID:      1,
		Name:    "Netflix",
		Amount:  amount,
		DateDue: time.Date(2020, time.December, 16, 0, 0, 0, 0, time.UTC),
	}

	user := userprofile.Userprofile{
		Name:  "Gary Gopher",
		Email: os.Getenv("EMAIL"),
	}

	t.Run("send an email", func(t *testing.T) {
		cal := calendar.CreateReminderInvite(subscription, reminder)
		client := &StubMailer{}
		datastore := &StubDataStore{subscription: subscription}

		err := SendEmail(reminder, user, cal, client, datastore)
		if err != nil {
			t.Errorf("there was an error sending the email %v", err)
		}

		if len(client.sentEmail.Attachments) != 1 {
			t.Errorf("no attachment recognised")
		}

		expectedSubject := fmt.Sprintf("Your %s subscription is due for renewal on %v", subscription.Name, subscription.DateDue.Format("January 2, 2006"))

		if client.sentEmail.Subject != expectedSubject {
			t.Errorf("did not get expected subject format, got %v want %v", client.sentEmail.Subject, expectedSubject)
		}
	})
}
