package calendar

import (
	"fmt"
	"time"

	"github.com/Catzkorn/subscrypt/internal/reminder"
	"github.com/Catzkorn/subscrypt/internal/subscription"
	ics "github.com/arran4/golang-ical"
)

const timeLayout = "January 2, 2006"

// CreateReminderInvite creates a new calendar invite
func CreateReminderInvite(subscription subscription.Subscription, reminder reminder.Reminder) *ics.Calendar {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	event := cal.AddEvent(fmt.Sprintf("%v@subscrypt.com", subscription.ID))
	event.SetCreatedTime(time.Now())
	event.SetDtStampTime(time.Now())
	event.SetModifiedAt(time.Now())
	event.SetAllDayStartAt(reminder.ReminderDate)
	event.SetSummary(fmt.Sprintf("Your %s subscription is due to renew on %v", subscription.Name, subscription.DateDue.Format(timeLayout)))
	event.SetLocation("")
	event.SetDescription(fmt.Sprintf("Hey! Your %s subscription is due to renew on %v and you asked us to remind you about that!",
		subscription.Name, subscription.DateDue))
	event.SetOrganizer("team@subscrypt.com", ics.WithCN("Subscrypt Team"))
	event.AddAttendee(fmt.Sprintf("%s", reminder.Email), ics.CalendarUserTypeIndividual, ics.ParticipationStatusNeedsAction, ics.ParticipationRoleReqParticipant, ics.WithRSVP(true))

	return cal

}
