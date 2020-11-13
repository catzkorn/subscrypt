package reminder

import "time"

// Reminder is the interface for reminder information
type Reminder struct {
	email          string
	subscriptionID int
	date           time.Time
}
