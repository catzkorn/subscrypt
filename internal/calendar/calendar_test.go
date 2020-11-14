package calendar

import (
	"testing"
	"time"
)

func TestCalendarInviteCreation(t *testing.T) {

	t.Run("checks that a .ics file is containg the correct information", func(t *testing.T) {

		_ = CreateReminderInvite(12, time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC), time.Date(2020, time.November, 16, 0, 0, 0, 0, time.UTC), "Netflix", "gary@gopher.com")

		// ioutil.WriteFile("./calendartest.ics", []byte(cal), 0644)

	})

}
