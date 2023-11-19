package alexa

import (
	"electricity-prices/pkg/i18n"
	"electricity-prices/pkg/price"
	"golang.org/x/text/language"
	"os"
	"strings"
	"testing"
)

var alexaService Service

func TestMain(m *testing.M) {

	i18n.InitialiseTranslations()
	alexaService = Service{}

	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestGetTodayNoDataMessage(t *testing.T) {
	testCases := []struct {
		name          string
		lang          language.Tag
		shouldContain string
	}{
		{
			name:          "English",
			lang:          language.English,
			shouldContain: "There is no data available yet for today",
		},
		{
			name:          "Spanish",
			lang:          language.Spanish,
			shouldContain: "No hay datos disponibles para hoy",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := alexaService.getTodayNoDataMessage(tc.lang)
			if !strings.Contains(actual, tc.shouldContain) {
				t.Errorf("shouldContain '%s' to contain: '%s'", actual, tc.shouldContain)
			}
		})
	}
}

func TestGetTomorrowNoDataMessage(t *testing.T) {
	testCases := []struct {
		name          string
		lang          language.Tag
		shouldContain string
	}{
		{
			name:          "English",
			lang:          language.English,
			shouldContain: "There is no data available yet for tomorrow",
		},
		{
			name:          "Spanish",
			lang:          language.Spanish,
			shouldContain: "Aún no hay datos disponibles para mañana",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := alexaService.getTomorrowNoDataMessage(tc.lang)
			if !strings.Contains(actual, tc.shouldContain) {
				t.Errorf("shouldContain '%s' to contain: '%s'", actual, tc.shouldContain)
			}
		})
	}
}

func TestGetTodayRatingMessage(t *testing.T) {
	testCases := []struct {
		name           string
		rating         price.DayRating
		dayAverage     float64
		lang           language.Tag
		shouldContain1 string
		shouldContain2 string
	}{
		{
			name:           "English good",
			rating:         price.Good,
			dayAverage:     0.1,
			lang:           language.English,
			shouldContain1: "good",
			shouldContain2: "10 cents per kilowatt-hour",
		},
		{
			name:           "English normal",
			rating:         price.Normal,
			dayAverage:     0.15,
			lang:           language.English,
			shouldContain1: "normal",
			shouldContain2: "15 cents per kilowatt-hour",
		},
		{
			name:           "English bad",
			rating:         price.Bad,
			dayAverage:     0.2,
			lang:           language.English,
			shouldContain1: "bad",
			shouldContain2: "20 cents per kilowatt-hour",
		},
		{
			name:           "Spanish good",
			rating:         price.Good,
			dayAverage:     0.1,
			lang:           language.Spanish,
			shouldContain1: "bueno",
			shouldContain2: "10 céntimos por kilovatio-hora",
		},
		{
			name:           "Spanish normal",
			rating:         price.Normal,
			dayAverage:     0.15,
			lang:           language.Spanish,
			shouldContain1: "normal",
			shouldContain2: "15 céntimos por kilovatio-hora",
		},
		{
			name:           "Spanish",
			rating:         price.Bad,
			dayAverage:     0.2,
			lang:           language.Spanish,
			shouldContain1: "malo",
			shouldContain2: "20 céntimos por kilovatio-hora",
		},
		{
			name:           "Nil rating",
			rating:         price.Nil,
			dayAverage:     0.2,
			lang:           language.Spanish,
			shouldContain1: "No hay datos disponibles para hoy",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := alexaService.getTodayRatingMessage(tc.rating, tc.dayAverage, tc.lang)
			if !strings.Contains(actual, tc.shouldContain1) {
				t.Errorf("shouldContain '%s' to contain: '%s'", actual, tc.shouldContain1)
			}
			if !strings.Contains(actual, tc.shouldContain2) {
				t.Errorf("shouldContain '%s' to contain: '%s'", actual, tc.shouldContain2)
			}
		})
	}
}

func TestGetTomorrowRatingMessage(t *testing.T) {
	testCases := []struct {
		name           string
		rating         price.DayRating
		dayAverage     float64
		lang           language.Tag
		shouldContain1 string
		shouldContain2 string
	}{
		{
			name:           "English good",
			rating:         price.Good,
			dayAverage:     0.1,
			lang:           language.English,
			shouldContain1: "good",
			shouldContain2: "10 cents per kilowatt-hour",
		},
		{
			name:           "English normal",
			rating:         price.Normal,
			dayAverage:     0.15,
			lang:           language.English,
			shouldContain1: "normal",
			shouldContain2: "15 cents per kilowatt-hour",
		},
		{
			name:           "English bad",
			rating:         price.Bad,
			dayAverage:     0.2,
			lang:           language.English,
			shouldContain1: "bad",
			shouldContain2: "20 cents per kilowatt-hour",
		},
		{
			name:           "Spanish good",
			rating:         price.Good,
			dayAverage:     0.1,
			lang:           language.Spanish,
			shouldContain1: "bueno",
			shouldContain2: "10 céntimos por kilovatio-hora",
		},
		{
			name:           "Spanish normal",
			rating:         price.Normal,
			dayAverage:     0.15,
			lang:           language.Spanish,
			shouldContain1: "normal",
			shouldContain2: "15 céntimos por kilovatio-hora",
		},
		{
			name:           "Spanish",
			rating:         price.Bad,
			dayAverage:     0.2,
			lang:           language.Spanish,
			shouldContain1: "malo",
			shouldContain2: "20 céntimos por kilovatio-hora",
		},
		{
			name:           "Nil rating",
			rating:         price.Nil,
			dayAverage:     0.2,
			lang:           language.Spanish,
			shouldContain1: "Aún no hay datos disponibles para mañana",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := alexaService.getTomorrowRatingMessage(tc.rating, tc.dayAverage, tc.lang)
			if !strings.Contains(actual, tc.shouldContain1) {
				t.Errorf("shouldContain '%s' to contain: '%s'", actual, tc.shouldContain1)
			}
			if !strings.Contains(actual, tc.shouldContain2) {
				t.Errorf("shouldContain '%s' to contain: '%s'", actual, tc.shouldContain2)
			}
		})
	}
}
