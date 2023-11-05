package utils

import (
	"log"
	"time"
)

func ParseStartAndEndTimes(date time.Time, numberOfDays int) (time.Time, time.Time) {
	xDaysAgo := date.AddDate(0, 0, -numberOfDays)
	startOfGivenDay := time.Date(xDaysAgo.Year(), xDaysAgo.Month(), xDaysAgo.Day(), 0, 0, 0, 0, xDaysAgo.Location())
	start := startOfGivenDay.Add(-time.Second) // subtract one second from the start of the given day

	nextDay := time.Date(date.Year(), date.Month(), date.Day()+1, 0, 0, 0, 0, date.Location())
	end := nextDay.Add(-time.Second)
	return start, end
}

func ParseDate(dateStr string) (time.Time, error) {
	location, err := time.LoadLocation("Europe/Madrid")
	if err != nil {
		log.Fatal(err)
	}

	// Parse the date string
	date, err := time.ParseInLocation("2006-01-02", dateStr, location)
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}

func ParseToLocalDay(date time.Time) string {
	location, err := time.LoadLocation("Europe/Madrid")
	if err != nil {
		log.Fatal(err)
	}

	return date.In(location).Format("2006-01-02")
}

func FormatTime(date time.Time) string {
	location, err := time.LoadLocation("Europe/Madrid")
	if err != nil {
		log.Fatal(err)
	}
	return date.In(location).Format("3 PM")
}

func SameHour(date1 time.Time, date2 time.Time) bool {
	location, err := time.LoadLocation("Europe/Madrid")
	if err != nil {
		log.Fatal(err)
	}
	locDate1 := date1.In(location)
	locDate2 := date2.In(location)
	return locDate1.Hour() == locDate2.Hour() && locDate1.Day() == locDate2.Day() && locDate1.Month() == locDate2.Month() && locDate1.Year() == locDate2.Year()
}

func StartOfDay(date time.Time) time.Time {
	location, err := time.LoadLocation("Europe/Madrid")
	if err != nil {
		log.Fatal(err)
	}
	localisedDate := date.In(location)
	return time.Date(localisedDate.Year(), localisedDate.Month(), localisedDate.Day(), 0, 0, 0, 0, localisedDate.Location())
}
