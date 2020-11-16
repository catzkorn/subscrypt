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

func TestSendingAnEmail(t *testing.T) {

	reminder := reminder.Reminder{
		Email:          os.Getenv("EMAIL"),
		SubscriptionID: 1,
		ReminderDate:   time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC),
	}
	amount, _ := decimal.NewFromString("8.00")
	subscription := subscription.Subscription{
		ID:      1,
		Name:    "Netflix",
		Amount:  amount,
		DateDue: time.Date(2020, time.November, 16, 0, 0, 0, 0, time.UTC),
	}

	t.Run("send an email", func(t *testing.T) {

		cal := calendar.CreateReminderInvite(subscription, reminder)
		client := &StubMailer{}
		err := SendEmail(reminder, cal, client)

		if err != nil {
			t.Errorf("there was an error sending the email %v", err)
		}

		if len(client.sentEmail.Attachments) != 1 {
			t.Errorf("no attachment recognised")
		}

		fmt.Println(client.sentEmail.Headers)

		// if client.sentEmail.SendAt

	})

}
