package reminder

import "time"

// Reminder is the interface for reminder information
type Reminder struct {
	Email          string
	SubscriptionID int
	ReminderDate   time.Time
}
