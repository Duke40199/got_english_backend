package utils

import (
	"fmt"
	"time"
)

// GetCurrentTime function is used to get the current time in milliseconds.
func GetCurrentEpochTimeInMiliseconds() int64 {
	var now = time.Now()
	ts := now.UnixNano() / 1000000
	return ts
}

func GetTimesByPeriod(period string) (time.Time, time.Time) {
	var startDate time.Time
	var endDate time.Time
	var timeNow = time.Now()
	switch period {
	case "weekly":
		{
			fmt.Println("================= weekly filter")
			endDate = time.Now()
			year, week := endDate.ISOWeek()
			//Get week start date
			startDate = GetWeekStart(year, week)
			break
		}
	case "daily":
		{
			fmt.Println("================= daily filter")
			startDate = time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 0, 0, 0, 0, time.Local)
			endDate = time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 23, 59, 59, 0, time.Local)
			break
		}
	case "monthly":
		{
			fmt.Println("================= monthly filter")
			startDate = time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 0, 0, 0, 0, time.Local)
			endDate = time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 23, 59, 59, 0, time.Local)
			break
		}
	}
	return startDate, endDate
}
func GetWeekStart(year, week int) time.Time {

	// Start from the middle of the year:
	t := time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC)

	// Roll back to Monday:
	if weekday := t.Weekday(); weekday == time.Sunday {
		t = t.AddDate(0, 0, -6)
	} else {
		t = t.AddDate(0, 0, -int(weekday)+1)
	}
	// Difference in weeks:
	_, w := t.ISOWeek()
	t = t.AddDate(0, 0, (week-w)*7)

	return t
}
