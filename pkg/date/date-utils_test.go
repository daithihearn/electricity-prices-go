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
		name          string
		date          string
		expected      time.Time
		errorExpected bool
	}{{
		name:          "Parse date",
		date:          "2023-01-02",
		expected:      time.Date(2023, 1, 2, 0, 0, 0, 0, location),
		errorExpected: false,
	},
		{
			name:          "Invalid date",
			date:          "2023-01-02-01",
			expected:      time.Time{},
			errorExpected: true,
		}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseDate(tc.date)
			if !result.Equal(tc.expected) {
				t.Errorf("Expected %v but was %v", tc.expected, result)
			}
			if tc.errorExpected && err == nil {
				t.Errorf("Expected error but was nil")
			}
		})
	}
}

func TestParseToLocalDay(t *testing.T) {
	location, _ := time.LoadLocation("Europe/Madrid")
	testCases := []struct {
		name     string
		date     time.Time
		expected string
	}{
		{
			name:     "Parse to local day",
			date:     time.Date(2023, 1, 2, 0, 0, 0, 0, location),
			expected: "2023-01-02",
		},
		{
			name:     "UTC different day",
			date:     time.Date(2023, 1, 1, 23, 59, 59, 0, time.UTC),
			expected: "2023-01-02",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ParseToLocalDay(tc.date)
			if result != tc.expected {
				t.Errorf("Expected %v but was %v", tc.expected, result)
			}
		})
	}
}

func TestFormatTime(t *testing.T) {
	location, _ := time.LoadLocation("Europe/Madrid")
	testCases := []struct {
		name     string
		date     time.Time
		expected string
	}{
		{
			name:     "Format time",
			date:     time.Date(2023, 1, 2, 0, 0, 0, 0, location),
			expected: "12 AM",
		},
		{
			name:     "Format time",
			date:     time.Date(2023, 1, 2, 12, 0, 0, 0, location),
			expected: "12 PM",
		},
		{
			name:     "Format time",
			date:     time.Date(2023, 1, 2, 23, 0, 0, 0, location),
			expected: "11 PM",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := FormatTime(tc.date)
			if result != tc.expected {
				t.Errorf("Expected %v but was %v", tc.expected, result)
			}
		})
	}
}

func TestSameHour(t *testing.T) {
	location, _ := time.LoadLocation("Europe/Madrid")
	testCases := []struct {
		name     string
		date1    time.Time
		date2    time.Time
		expected bool
	}{
		{
			name:     "Same hour",
			date1:    time.Date(2023, 1, 2, 0, 0, 0, 0, location),
			date2:    time.Date(2023, 1, 2, 0, 0, 0, 0, location),
			expected: true,
		},
		{
			name:     "Same hour",
			date1:    time.Date(2023, 1, 2, 0, 0, 0, 0, location),
			date2:    time.Date(2023, 1, 2, 1, 0, 0, 0, location),
			expected: false,
		},
		{
			name:     "Same hour",
			date1:    time.Date(2023, 1, 2, 0, 0, 0, 0, location),
			date2:    time.Date(2023, 1, 3, 0, 0, 0, 0, location),
			expected: false,
		},
		{
			name:     "Same hour",
			date1:    time.Date(2023, 1, 2, 0, 0, 0, 0, location),
			date2:    time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := SameHour(tc.date1, tc.date2)
			if result != tc.expected {
				t.Errorf("Expected %v but was %v", tc.expected, result)
			}
		})
	}
}

func TestStartOfDay(t *testing.T) {
	location, _ := time.LoadLocation("Europe/Madrid")
	testCases := []struct {
		name     string
		date     time.Time
		expected time.Time
	}{
		{
			name:     "Start of day",
			date:     time.Date(2023, 1, 2, 0, 0, 0, 0, location),
			expected: time.Date(2023, 1, 2, 0, 0, 0, 0, location),
		},
		{
			name:     "Start of day",
			date:     time.Date(2023, 1, 2, 23, 59, 59, 0, location),
			expected: time.Date(2023, 1, 2, 0, 0, 0, 0, location),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := StartOfDay(tc.date)
			if !result.Equal(tc.expected) {
				t.Errorf("Expected %v but was %v", tc.expected, result)
			}
		})
	}
}

func TestParseEsiosTime(t *testing.T) {
	location, _ := time.LoadLocation("Europe/Madrid")
	testCases := []struct {
		name          string
		date          string
		hourRange     string
		expected      time.Time
		errorExpected bool
	}{
		{
			name:          "Parse esios time - start of day",
			date:          "02/01/2023",
			hourRange:     "00-01",
			expected:      time.Date(2023, 1, 2, 0, 0, 0, 0, location),
			errorExpected: false,
		},
		{
			name:          "Parse esios time - middle of day",
			date:          "02/01/2023",
			hourRange:     "12-13",
			expected:      time.Date(2023, 1, 2, 12, 0, 0, 0, location),
			errorExpected: false,
		},
		{
			name:          "Parse esios time - end of day",
			date:          "02/01/2023",
			hourRange:     "23-24",
			expected:      time.Date(2023, 1, 2, 23, 0, 0, 0, location),
			errorExpected: false,
		},
		{
			name:          "Parse esios time - invalid date",
			date:          "invalid date",
			hourRange:     "23-24",
			expected:      time.Time{},
			errorExpected: true,
		},
		{
			name:          "Parse esios time - invalid hour range",
			date:          "02/01/2023",
			hourRange:     "invalid hour range",
			expected:      time.Time{},
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseEsiosTime(tc.date, tc.hourRange)
			if !result.Equal(tc.expected) {
				t.Errorf("Expected %v but was %v", tc.expected, result)
			}
			if tc.errorExpected && err == nil {
				t.Errorf("Expected error but was nil")
			}
		})
	}
}
