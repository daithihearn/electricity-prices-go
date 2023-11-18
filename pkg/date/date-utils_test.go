package date

import (
	"testing"
	"time"
)

func TestParseStartAndEndTimes(t *testing.T) {
	startOfADay := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endOfADay := time.Date(2023, 1, 1, 23, 59, 59, 0, time.UTC)
	location, _ := time.LoadLocation("Europe/Madrid")
	startOfADayInMadrid := time.Date(2023, 1, 1, 0, 0, 0, 0, location)
	endOfADayInMadrid := time.Date(2023, 1, 1, 23, 59, 59, 0, location)

	testCases := []struct {
		name          string
		endDate       time.Time
		numberOfDays  int
		expectedStart time.Time
		expectedEnd   time.Time
	}{
		{
			name:          "Start of Day UTC",
			endDate:       startOfADay,
			numberOfDays:  20,
			expectedStart: time.Date(2022, 12, 12, 23, 59, 59, 0, time.UTC),
			expectedEnd:   time.Date(2023, 1, 1, 23, 59, 59, 0, time.UTC),
		},
		{
			name:          "End of Day UTC",
			endDate:       endOfADay,
			numberOfDays:  20,
			expectedStart: time.Date(2022, 12, 12, 23, 59, 59, 0, time.UTC),
			expectedEnd:   time.Date(2023, 1, 1, 23, 59, 59, 0, time.UTC),
		},
		{
			name:          "Start of Day Madrid",
			endDate:       startOfADayInMadrid,
			numberOfDays:  10,
			expectedStart: time.Date(2022, 12, 22, 23, 59, 59, 0, location),
			expectedEnd:   time.Date(2023, 1, 1, 23, 59, 59, 0, location),
		},
		{
			name:          "End of Day Madrid",
			endDate:       endOfADayInMadrid,
			numberOfDays:  10,
			expectedStart: time.Date(2022, 12, 22, 23, 59, 59, 0, location),
			expectedEnd:   time.Date(2023, 1, 1, 23, 59, 59, 0, location),
		},
		{
			name:          "Only one day",
			endDate:       time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			numberOfDays:  1,
			expectedStart: time.Date(2023, 1, 1, 23, 59, 59, 0, time.UTC),
			expectedEnd:   time.Date(2023, 1, 2, 23, 59, 59, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			start, end := ParseStartAndEndTimes(tc.endDate, tc.numberOfDays)
			if start != tc.expectedStart {
				t.Errorf("Expected start to be %v but was %v", tc.expectedStart, start)
			}
			if end != tc.expectedEnd {
				t.Errorf("Expected end to be %v but was %v", tc.expectedEnd, end)
			}
		})
	}

}

func TestParseDate(t *testing.T) {
	location, _ := time.LoadLocation("Europe/Madrid")
	testCases := []struct {
		name     string
		date     string
		expected time.Time
	}{{
		name:     "Parse date",
		date:     "2023-01-02",
		expected: time.Date(2023, 1, 2, 0, 0, 0, 0, location),
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, _ := ParseDate(tc.date)
			if !result.Equal(tc.expected) {
				t.Errorf("Expected %v but was %v", tc.expected, result)
			}
		})
	}
}
