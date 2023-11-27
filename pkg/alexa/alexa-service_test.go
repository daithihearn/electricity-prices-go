package alexa

import (
	"electricity-prices/pkg/i18n"
	"electricity-prices/pkg/price"
	"golang.org/x/text/language"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

var alexaService Service
var madridLocation *time.Location

var period1 []price.Price
var period2 []price.Price
var period3 []price.Price
var period4 []price.Price

func TestMain(m *testing.M) {

	i18n.InitialiseTranslations()
	alexaService = Service{}

	var err error
	madridLocation, err = time.LoadLocation("Europe/Madrid")
	if err != nil {
		log.Fatal(err)
	}

	period1 = []price.Price{
		{
			DateTime: time.Date(2023, 1, 1, 3, 0, 0, 0, madridLocation),
			Price:    0.1,
		},
		{
			DateTime: time.Date(2023, 1, 1, 4, 0, 0, 0, madridLocation),
			Price:    0.2,
		},
	}

	period2 = []price.Price{
		{
			DateTime: time.Date(2023, 1, 1, 6, 0, 0, 0, madridLocation),
			Price:    0.3,
		},
		{
			DateTime: time.Date(2023, 1, 1, 7, 0, 0, 0, madridLocation),
			Price:    0.4,
		},
	}

	period3 = []price.Price{
		{
			DateTime: time.Date(2023, 1, 1, 9, 0, 0, 0, madridLocation),
			Price:    0.3,
		},
		{
			DateTime: time.Date(2023, 1, 1, 10, 0, 0, 0, madridLocation),
			Price:    0.4,
		},
	}

	period4 = []price.Price{
		{
			DateTime: time.Date(2023, 1, 1, 16, 0, 0, 0, madridLocation),
			Price:    0.3,
		},
	}

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

func TestGetNextExpensivePeriodMessage(t *testing.T) {

	testCases := []struct {
		name           string
		t              time.Time
		periods        [][]price.Price
		lang           language.Tag
		shouldContain1 string
		shouldContain2 string
		shouldContain3 string
	}{
		{
			name: "Period in the future (English)",
			t:    time.Date(2023, 1, 1, 1, 0, 0, 0, madridLocation),
			periods: [][]price.Price{
				period1, period2,
			},
			lang:           language.English,
			shouldContain1: "next expensive period starts at 3 AM",
			shouldContain2: "15 cents per kilowatt-hour",
			shouldContain3: "end at 5 AM",
		},
		{
			name: "Period in the future (Spanish)",
			t:    time.Date(2023, 1, 1, 1, 0, 0, 0, madridLocation),
			periods: [][]price.Price{
				period3, period4,
			},
			lang:           language.Spanish,
			shouldContain1: "El próximo período caro comienza a las 9 AM",
			shouldContain2: "35 céntimos por kilovatio-hora",
			shouldContain3: "terminara a las 11 AM",
		},
		{
			name: "Period has started (English)",
			t:    time.Date(2023, 1, 1, 3, 30, 0, 0, madridLocation),
			periods: [][]price.Price{
				period1, period2, period3,
			},
			lang:           language.English,
			shouldContain1: "You are currently in an expensive period that started at 3 AM",
			shouldContain2: "15 cents per kilowatt-hour",
			shouldContain3: "end at 5 AM",
		},
		{
			name: "Period has started (Spanish)",
			t:    time.Date(2023, 1, 1, 9, 30, 0, 0, madridLocation),
			periods: [][]price.Price{
				period3, period4,
			},
			lang:           language.Spanish,
			shouldContain1: "Actualmente se encuentra en un período caro que comenzó a las 9 AM",
			shouldContain2: "35 céntimos por kilovatio-hora",
			shouldContain3: "terminara a las 11 AM",
		},
		{
			name:           "No data (English)",
			t:              time.Date(2023, 1, 1, 1, 0, 0, 0, madridLocation),
			periods:        [][]price.Price{},
			lang:           language.English,
			shouldContain1: "There are no expensive periods today.",
			shouldContain2: "",
			shouldContain3: "",
		},
		{
			name:           "No data (Spanish)",
			t:              time.Date(2023, 1, 1, 1, 0, 0, 0, madridLocation),
			periods:        [][]price.Price{},
			lang:           language.Spanish,
			shouldContain1: "Hoy no hay ningun período caro.",
			shouldContain2: "",
			shouldContain3: "",
		},
		{
			name: "Periods have ended (English)",
			t:    time.Date(2023, 1, 1, 8, 0, 0, 0, madridLocation),
			periods: [][]price.Price{
				period1, period2,
			},
			lang:           language.English,
			shouldContain1: "The expensive periods for today have already passed.",
			shouldContain2: "",
			shouldContain3: "",
		},
		{
			name: "Periods have ended (Spanish)",
			t:    time.Date(2023, 1, 1, 8, 0, 0, 0, madridLocation),
			periods: [][]price.Price{
				period1, period2,
			},
			lang:           language.Spanish,
			shouldContain1: "Los períodos caros de hoy ya ha pasado.",
			shouldContain2: "",
			shouldContain3: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := alexaService.getNextExpensivePeriodMessage(tc.periods, tc.t, tc.lang)
			if !strings.Contains(actual, tc.shouldContain1) {
				t.Errorf("'%s' should contain: '%s'", actual, tc.shouldContain1)
			}
			if !strings.Contains(actual, tc.shouldContain2) {
				t.Errorf("'%s' should contain: '%s'", actual, tc.shouldContain2)
			}
			if !strings.Contains(actual, tc.shouldContain3) {
				t.Errorf("'%s' should contain: '%s'", actual, tc.shouldContain3)
			}
		})
	}
}

func TestGetNextCheapPeriodMessage(t *testing.T) {

	testCases := []struct {
		name           string
		t              time.Time
		periods        [][]price.Price
		lang           language.Tag
		shouldContain1 string
		shouldContain2 string
		shouldContain3 string
	}{
		{
			name: "Period in the future (English)",
			t:    time.Date(2023, 1, 1, 1, 0, 0, 0, madridLocation),
			periods: [][]price.Price{
				period1, period2,
			},
			lang:           language.English,
			shouldContain1: "next cheap period starts at 3 AM",
			shouldContain2: "15 cents per kilowatt-hour",
			shouldContain3: "end at 5 AM",
		},
		{
			name: "Period in the future (Spanish)",
			t:    time.Date(2023, 1, 1, 1, 0, 0, 0, madridLocation),
			periods: [][]price.Price{
				period3, period4,
			},
			lang:           language.Spanish,
			shouldContain1: "El próximo período barato comienza a las 9 AM",
			shouldContain2: "35 céntimos por kilovatio-hora",
			shouldContain3: "terminara a las 11 AM",
		},
		{
			name: "Period has started (English)",
			t:    time.Date(2023, 1, 1, 3, 30, 0, 0, madridLocation),
			periods: [][]price.Price{
				period1, period2, period3,
			},
			lang:           language.English,
			shouldContain1: "You are currently in a cheap period that started at 3 AM",
			shouldContain2: "15 cents per kilowatt-hour",
			shouldContain3: "end at 5 AM",
		},
		{
			name: "Period has started (Spanish)",
			t:    time.Date(2023, 1, 1, 9, 30, 0, 0, madridLocation),
			periods: [][]price.Price{
				period3, period4,
			},
			lang:           language.Spanish,
			shouldContain1: "Actualmente se encuentra en un período barato que comenzó a las 9 AM",
			shouldContain2: "35 céntimos por kilovatio-hora",
			shouldContain3: "terminara a las 11 AM",
		},
		{
			name:           "No data (English)",
			t:              time.Date(2023, 1, 1, 1, 0, 0, 0, madridLocation),
			periods:        [][]price.Price{},
			lang:           language.English,
			shouldContain1: "There are no cheap periods today.",
			shouldContain2: "",
			shouldContain3: "",
		},
		{
			name:           "No data (Spanish)",
			t:              time.Date(2023, 1, 1, 1, 0, 0, 0, madridLocation),
			periods:        [][]price.Price{},
			lang:           language.Spanish,
			shouldContain1: "Hoy no hay ningun período barato.",
			shouldContain2: "",
			shouldContain3: "",
		},
		{
			name: "Periods have ended (English)",
			t:    time.Date(2023, 1, 1, 8, 0, 0, 0, madridLocation),
			periods: [][]price.Price{
				period1, period2,
			},
			lang:           language.English,
			shouldContain1: "The cheap periods for today have already passed.",
			shouldContain2: "",
			shouldContain3: "",
		},
		{
			name: "Periods have ended (Spanish)",
			t:    time.Date(2023, 1, 1, 8, 0, 0, 0, madridLocation),
			periods: [][]price.Price{
				period1, period2,
			},
			lang:           language.Spanish,
			shouldContain1: "Los periodos baratos de hoy ya han pasado.",
			shouldContain2: "",
			shouldContain3: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := alexaService.getNextCheapPeriodMessage(tc.periods, tc.t, tc.lang)
			if !strings.Contains(actual, tc.shouldContain1) {
				t.Errorf("'%s' should contain: '%s'", actual, tc.shouldContain1)
			}
			if !strings.Contains(actual, tc.shouldContain2) {
				t.Errorf("'%s' should contain: '%s'", actual, tc.shouldContain2)
			}
			if !strings.Contains(actual, tc.shouldContain3) {
				t.Errorf("'%s' should contain: '%s'", actual, tc.shouldContain3)
			}
		})
	}
}
