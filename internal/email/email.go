package email

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/Catzkorn/subscrypt/internal/reminder"
	ics "github.com/arran4/golang-ical"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendEmail sends a reminder email
func SendEmail(reminder reminder.Reminder, event *ics.Calendar) error {
	email := reminder.Email
	subscription := reminder.SubscriptionID
	reminderDate := reminder.ReminderDate

	from := mail.NewEmail("Go, Team!", "reminder@subscrypt.com")
	subject := fmt.Sprintf("Your %d subscription is due for renewal on %v", subscription, reminderDate)
	to := mail.NewEmail("Hacker", email)
	plainTextContent := "Soon we will have calendar functionality!"
	htmlContent := "<strong>Soon we will have calendar functionality!</strong>"

	calendarInvite := mail.NewAttachment()

	encoded := base64.StdEncoding.EncodeToString([]byte(event.Serialize()))
	calendarInvite.SetContent(encoded)
	calendarInvite.SetType("text/plain")
	calendarInvite.SetFilename("subscriptionreminder.ics")
	calendarInvite.SetDisposition("attachment")

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	message.AddAttachment(calendarInvite)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		return err
	}
	fmt.Println(response.StatusCode)
	fmt.Println(response.Body)
	fmt.Println(response.Headers)
	return nil
}
