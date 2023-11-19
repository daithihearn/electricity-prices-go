package date

import (
	"log"
	"time"
)

var madridLocation *time.Location

func init() {
	var err error
	madridLocation, err = time.LoadLocation("Europe/Madrid")
	if err != nil {
		log.Fatal(err)
	}
}

// ParseStartAndEndTimes
// Given an end date and a number of days, calculate the start and end dates for use in a query.
func ParseStartAndEndTimes(endDate time.Time, numberOfDays int) (time.Time, time.Time) {
	xDaysAgo := endDate.AddDate(0, 0, -numberOfDays)
	start := time.Date(xDaysAgo.Year(), xDaysAgo.Month(), xDaysAgo.Day(), 23, 59, 59, 0, xDaysAgo.Location())

	nextDay := time.Date(endDate.Year(), endDate.Month(), endDate.Day()+1, 0, 0, 0, 0, endDate.Location())
	end := nextDay.Add(-time.Second)
	return start, end
}

func ParseDate(dateStr string) (time.Time, error) {
	// Parse the date string
	date, err := time.ParseInLocation("2006-01-02", dateStr, madridLocation)
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}

func ParseToLocalDay(date time.Time) string {
	return date.In(madridLocation).Format("2006-01-02")
}

func FormatTime(date time.Time) string {
	return date.In(madridLocation).Format("3 PM")
}

func SameHour(date1 time.Time, date2 time.Time) bool {
	locDate1 := date1.In(madridLocation)
	locDate2 := date2.In(madridLocation)
	return locDate1.Hour() == locDate2.Hour() && locDate1.Day() == locDate2.Day() && locDate1.Month() == locDate2.Month() && locDate1.Year() == locDate2.Year()
}

func StartOfDay(date time.Time) time.Time {
	localisedDate := date.In(madridLocation)
	return time.Date(localisedDate.Year(), localisedDate.Month(), localisedDate.Day(), 0, 0, 0, 0, localisedDate.Location())
}
