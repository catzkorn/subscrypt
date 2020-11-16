package email

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/Catzkorn/subscrypt/internal/reminder"
	"github.com/Catzkorn/subscrypt/internal/subscription"
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
func SendEmail(reminder reminder.Reminder, event *ics.Calendar, mailer Mailer, datastore DataStore) error {

	subscription, err := datastore.GetSubscription(reminder.SubscriptionID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	from := mail.NewEmail("Subscrypt Team", "team@subscrypt.com")
	subject := fmt.Sprintf("Your %s subscription is due for renewal on %v", subscription.Name, reminder.ReminderDate.Format(timeLayout))
	to := mail.NewEmail("Subscryptee", reminder.Email)
	plainTextContent := "Hey there Subscryptee!\nYou asked for a reminder and here it is!"
	htmlContent := "<strong>Hey there Subscryptee!\nYou asked for a reminder and here it is!</strong>"

	calendarInvite := createAttachment(event)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	message.AddAttachment(calendarInvite)

	// client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
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
	calendarInvite.SetFilename("subscriptionreminder.ics")
	calendarInvite.SetDisposition("attachment")

	return calendarInvite
}
