package email

import (
	"fmt"
	"log"
	"os"

	"github.com/Catzkorn/subscrypt/internal/reminder"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendEmail sends a reminder email
func SendEmail(reminder reminder.Reminder) {
	email := reminder.Email
	subscription := reminder.SubscriptionID
	reminderDate := reminder.ReminderDate

	from := mail.NewEmail("Go, Team!", "reminder@subscrypt.com")
	subject := fmt.Sprintf("Your %d subscription is due for renewal on %v", subscription, reminderDate)
	to := mail.NewEmail("Hacker", email)
	plainTextContent := "Soon we will have calendar functionality!"
	htmlContent := "<strong>Soon we will have calendar functionality!</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
