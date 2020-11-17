package email

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/Catzkorn/subscrypt/internal/reminder"
	"github.com/Catzkorn/subscrypt/internal/subscription"
	"github.com/Catzkorn/subscrypt/internal/userprofile"
	ics "github.com/arran4/golang-ical"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Mailer defines the interface required to send an email
type Mailer interface {
	Send(email *mail.SGMailV3) (*rest.Response, error)
}

// DataStore defines the interface required to get a subscription
type DataStore interface {
	GetSubscription(subscriptionID int) (*subscription.Subscription, error)
}

const timeLayout = "January 2, 2006"

// SendEmail sends a reminder email
func SendEmail(reminder reminder.Reminder, user userprofile.Userprofile, event *ics.Calendar, mailer Mailer, datastore DataStore) error {

	subscription, err := datastore.GetSubscription(reminder.SubscriptionID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}
	if subscription == nil {
		return fmt.Errorf("no subscription found for ID: %v", reminder.SubscriptionID)
	}

	from := mail.NewEmail("Subscrypt Team", "team@subscrypt.com")
	subject := fmt.Sprintf("Your %s subscription is due for renewal on %v", subscription.Name, subscription.DateDue.Format(timeLayout))
	to := mail.NewEmail(user.Name, reminder.Email)
	plainTextContent := fmt.Sprintf("Hey there %s!\nYou asked for a reminder and here it is!", user.Name)
	htmlContent := fmt.Sprintf("<strong>Hey there %s!\nYou asked for a reminder and here it is!</strong>", user.Name)

	calendarInvite := createAttachment(event)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	message.AddAttachment(calendarInvite)

	response, err := mailer.Send(message)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("did not return expected status code")
	}

	return nil
}

func createAttachment(event *ics.Calendar) *mail.Attachment {
	calendarInvite := mail.NewAttachment()

	encoded := base64.StdEncoding.EncodeToString([]byte(event.Serialize()))
	calendarInvite.SetContent(encoded)
	calendarInvite.SetType("text/plain")
	calendarInvite.SetFilename("subscryptreminder.ics")
	calendarInvite.SetDisposition("attachment")

	return calendarInvite
}
