package calendar

import (
	"fmt"
	"time"

	ics "github.com/arran4/golang-ical"
)

// CreateReminderInvite creates a new calendar invite
func CreateReminderInvite(subscriptionID int, reminderDate time.Time, subscriptionDate time.Time, subscriptionName string, email string) *ics.Calendar {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	event := cal.AddEvent(fmt.Sprintf("%v@subscrypt.com", subscriptionID))
	event.SetCreatedTime(time.Now())
	event.SetDtStampTime(time.Now())
	event.SetModifiedAt(time.Now())
	event.SetAllDayStartAt(reminderDate)
	event.SetSummary(fmt.Sprintf("Your %s subscription is due to renew on %v", "Netflix", subscriptionDate))
	event.SetLocation("")
	event.SetDescription(fmt.Sprintf("Hey! Your %s subscription is due to renew on %v and you asked us to remind you about that!",
		subscriptionName, subscriptionDate))
	event.SetOrganizer("team@subscrypt.com", ics.WithCN("Subscrypt Team"))
	event.AddAttendee(fmt.Sprintf("%s", email), ics.CalendarUserTypeIndividual, ics.ParticipationStatusNeedsAction, ics.ParticipationRoleReqParticipant, ics.WithRSVP(true))

	return cal

}
