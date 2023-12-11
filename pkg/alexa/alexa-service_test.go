package alexa

import (
	"context"
	"electricity-prices/pkg/i18n"
	"electricity-prices/pkg/price"
	"errors"
	"golang.org/x/text/language"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

var alexaService Service
var madridLocation *time.Location

var pricesToday []price.Price

var pricesTomorrow []price.Price

var period1 []price.Price
var period2 []price.Price
var period3 []price.Price
var period4 []price.Price

func TestMain(m *testing.M) {

	err := i18n.InitialiseTranslations(
		[]i18n.File{
			{
				Filename: "en.toml",
				Lang:     language.English,
			},
			{
				Filename: "es.toml",
				Lang:     language.Spanish,
			},
		})
	if err != nil {
		log.Fatal(err)
	}
	alexaService = Service{}

	madridLocation, err = time.LoadLocation("Europe/Madrid")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 24; i++ {
		pricesToday = append(pricesToday, price.Price{
			DateTime: time.Date(2023, 1, 1, i, 0, 0, 0, madridLocation),
			Price:    0.01 + float64(i)/100,
		})

		pricesTomorrow = append(pricesTomorrow, price.Price{
			DateTime: time.Date(2023, 1, 2, i, 0, 0, 0, madridLocation),
			Price:    0.02 + float64(i)/100,
		})
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

func TestGetTitle(t *testing.T) {
	testCases := []struct {
		name          string
		lang          language.Tag
		shouldContain string
	}{
		{
			name:          "English",
			lang:          language.English,
			shouldContain: "Electricity Prices",
		},
		{
			name:          "Spanish",
			lang:          language.Spanish,
			shouldContain: "Precio de la electricidad",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := alexaService.GetTitle(tc.lang)
			if !strings.Contains(actual, tc.shouldContain) {
				t.Errorf("Expected '%s' to contain: '%s'", actual, tc.shouldContain)
			}
		})
	}
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
			shouldContain1: "Los períodos caros de hoy ya han pasado.",
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

func TestGetPriceMessage(t *testing.T) {
	testCases := []struct {
		name           string
		prices         []price.Price
		t              time.Time
		lang           language.Tag
		shouldContain1 string
	}{
		{
			name: "Current price (English)",
			prices: []price.Price{
				{
					DateTime: time.Date(2023, 1, 1, 3, 0, 0, 0, madridLocation),
					Price:    0.1,
				},
				{
					DateTime: time.Date(2023, 1, 1, 4, 0, 0, 0, madridLocation),
					Price:    0.2,
				},
			},
			t:              time.Date(2023, 1, 1, 3, 30, 0, 0, madridLocation),
			lang:           language.English,
			shouldContain1: "The current price is 10 cents per kilowatt-hour",
		},
		{
			name: "Current price (Spanish)",
			prices: []price.Price{
				{
					DateTime: time.Date(2023, 1, 1, 3, 0, 0, 0, madridLocation),
					Price:    0.1,
				},
				{
					DateTime: time.Date(2023, 1, 1, 4, 0, 0, 0, madridLocation),
					Price:    0.2,
				},
			},
			t:              time.Date(2023, 1, 1, 4, 30, 0, 0, madridLocation),
			lang:           language.Spanish,
			shouldContain1: "El precio actual es 20 céntimos por kilovatio-hora",
		},
		{
			name: "No data (English)",
			prices: []price.Price{
				{
					DateTime: time.Date(2023, 1, 1, 3, 0, 0, 0, madridLocation),
					Price:    0.1,
				},
				{
					DateTime: time.Date(2023, 1, 1, 4, 0, 0, 0, madridLocation),
					Price:    0.2,
				},
			},
			t:              time.Date(2023, 1, 1, 5, 30, 0, 0, madridLocation),
			lang:           language.English,
			shouldContain1: "There is no data available yet for today",
		},
		{
			name: "No data (Spanish)",
			prices: []price.Price{
				{
					DateTime: time.Date(2023, 1, 1, 3, 0, 0, 0, madridLocation),
					Price:    0.1,
				},
				{
					DateTime: time.Date(2023, 1, 1, 4, 0, 0, 0, madridLocation),
					Price:    0.2,
				},
			},
			t:              time.Date(2023, 1, 1, 5, 30, 0, 0, madridLocation),
			lang:           language.Spanish,
			shouldContain1: "No hay datos disponibles para hoy",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := alexaService.getPriceMessage(tc.prices, tc.t, tc.lang)
			if !strings.Contains(actual, tc.shouldContain1) {
				t.Errorf("'%s' should contain: '%s'", actual, tc.shouldContain1)
			}
		})
	}
}

func TestGetFullFeed(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name            string
		t               time.Time
		mockResultToday price.DailyPriceInfo
		mockErrorToday  error
		mockResultTmrw  price.DailyPriceInfo
		mockErrorTmrw   error
		lang            language.Tag
		shouldContain1  string
		shouldContain2  string
		shouldContain3  string
		shouldContain4  string
		shouldContain5  string
	}{
		{
			name: "Full feed - start of day (English)",
			t:    time.Date(2023, 1, 1, 1, 0, 0, 0, madridLocation),
			mockResultToday: price.DailyPriceInfo{
				DayRating:        price.Good,
				DayAverage:       0.1,
				ThirtyDayAverage: 0.2,
				Prices:           pricesToday,
				ExpensivePeriods: [][]price.Price{
					pricesToday[2:5], pricesToday[8:9],
				},
				CheapPeriods: [][]price.Price{
					pricesToday[10:11], pricesToday[13:14],
				},
			},
			mockErrorToday: nil,
			mockResultTmrw: price.DailyPriceInfo{
				DayRating:        price.Normal,
				DayAverage:       0.15,
				ThirtyDayAverage: 0.25,
				Prices:           pricesTomorrow,
				ExpensivePeriods: [][]price.Price{
					pricesTomorrow[2:5], pricesTomorrow[8:9],
				},
				CheapPeriods: [][]price.Price{
					pricesTomorrow[4:5],
				},
			},
			mockErrorTmrw:  nil,
			lang:           language.English,
			shouldContain1: "Today is a good day with an average price of 10 cents per kilowatt-hour",
			shouldContain2: "The current price is 2 cents per kilowatt-hour",
			shouldContain3: "The next cheap period starts at 10 AM with an average price of 11 cents per kilowatt-hour and will end at 11 AM",
			shouldContain4: "The next expensive period starts at 2 AM with an average price of 4 cents per kilowatt-hour and will end at 5 AM",
			shouldContain5: "Tomorrow is a normal day with an average price of 15 cents per kilowatt-hour",
		},
		{
			name: "Full feed - start of day (Spanish)",
			t:    time.Date(2023, 1, 1, 1, 0, 0, 0, madridLocation),
			mockResultToday: price.DailyPriceInfo{
				DayRating:        price.Good,
				DayAverage:       0.1,
				ThirtyDayAverage: 0.2,
				Prices:           pricesToday,
				ExpensivePeriods: [][]price.Price{
					pricesToday[2:5], pricesToday[8:9],
				},
				CheapPeriods: [][]price.Price{
					pricesToday[10:11], pricesToday[13:14],
				},
			},
			mockErrorToday: nil,
			mockResultTmrw: price.DailyPriceInfo{
				DayRating:        price.Normal,
				DayAverage:       0.15,
				ThirtyDayAverage: 0.25,
				Prices:           pricesTomorrow,
				ExpensivePeriods: [][]price.Price{
					pricesTomorrow[2:5], pricesTomorrow[8:9],
				},
				CheapPeriods: [][]price.Price{
					pricesTomorrow[4:5],
				},
			},
			mockErrorTmrw:  nil,
			lang:           language.Spanish,
			shouldContain1: "Hoy es un día bueno, con un precio medio de 10 céntimos por kilovatio-hora.",
			shouldContain2: "El precio actual es 2 céntimos por kilovatio-hora.",
			shouldContain3: "El próximo período barato comienza a las 10 AM con un precio medio de 11 céntimos por kilovatio-hora y terminara a las 11 AM.",
			shouldContain4: "El próximo período caro comienza a las 2 AM con un precio promedio de 4 céntimos por kilovatio-hora y terminara a las 5 AM.",
			shouldContain5: "Mañana es un día normal, con un precio promedio de 15 céntimos por kilovatio-hora.",
		},
		{
			name: "Full feed - first cheap period passed",
			t:    time.Date(2023, 1, 1, 12, 0, 0, 0, madridLocation),
			mockResultToday: price.DailyPriceInfo{
				DayRating:        price.Good,
				DayAverage:       0.1,
				ThirtyDayAverage: 0.2,
				Prices:           pricesToday,
				ExpensivePeriods: [][]price.Price{
					pricesToday[2:5], pricesToday[8:9],
				},
				CheapPeriods: [][]price.Price{
					pricesToday[10:11], pricesToday[13:14],
				},
			},
			mockErrorToday: nil,
			mockResultTmrw: price.DailyPriceInfo{
				DayRating:        price.Normal,
				DayAverage:       0.15,
				ThirtyDayAverage: 0.25,
				Prices:           pricesTomorrow,
				ExpensivePeriods: [][]price.Price{
					pricesTomorrow[2:5], pricesTomorrow[8:9],
				},
				CheapPeriods: [][]price.Price{
					pricesTomorrow[4:5],
				},
			},
			mockErrorTmrw:  nil,
			lang:           language.English,
			shouldContain1: "Today is a good day with an average price of 10 cents per kilowatt-hour",
			shouldContain2: "The current price is 13 cents per kilowatt-hour",
			shouldContain3: "The next cheap period starts at 1 PM with an average price of 14 cents per kilowatt-hour and will end at 2 PM",
			shouldContain4: "The expensive periods for today have already passed.",
			shouldContain5: "Tomorrow is a normal day with an average price of 15 cents per kilowatt-hour",
		},
		{
			name: "Full feed - first expensive period passed",
			t:    time.Date(2023, 1, 1, 6, 0, 0, 0, madridLocation),
			mockResultToday: price.DailyPriceInfo{
				DayRating:        price.Good,
				DayAverage:       0.1,
				ThirtyDayAverage: 0.2,
				Prices:           pricesToday,
				ExpensivePeriods: [][]price.Price{
					pricesToday[2:5], pricesToday[8:9],
				},
				CheapPeriods: [][]price.Price{
					pricesToday[10:11], pricesToday[13:14],
				},
			},
			mockErrorToday: nil,
			mockResultTmrw: price.DailyPriceInfo{
				DayRating:        price.Normal,
				DayAverage:       0.15,
				ThirtyDayAverage: 0.25,
				Prices:           pricesTomorrow,
				ExpensivePeriods: [][]price.Price{
					pricesTomorrow[2:5], pricesTomorrow[8:9],
				},
				CheapPeriods: [][]price.Price{
					pricesTomorrow[4:5],
				},
			},
			mockErrorTmrw:  nil,
			lang:           language.English,
			shouldContain1: "Today is a good day with an average price of 10 cents per kilowatt-hour",
			shouldContain2: "The current price is 7 cents per kilowatt-hour",
			shouldContain3: "The next cheap period starts at 10 AM with an average price of 11 cents per kilowatt-hour and will end at 11 AM",
			shouldContain4: "The next expensive period starts at 8 AM with an average price of 9 cents per kilowatt-hour and will end at 9 AM",
			shouldContain5: "Tomorrow is a normal day with an average price of 15 cents per kilowatt-hour",
		},
		{
			name: "Full feed - end of day",
			t:    time.Date(2023, 1, 1, 23, 0, 0, 0, madridLocation),
			mockResultToday: price.DailyPriceInfo{
				DayRating:        price.Good,
				DayAverage:       0.1,
				ThirtyDayAverage: 0.2,
				Prices:           pricesToday,
				ExpensivePeriods: [][]price.Price{
					pricesToday[2:5], pricesToday[8:9],
				},
				CheapPeriods: [][]price.Price{
					pricesToday[10:11], pricesToday[13:14],
				},
			},
			mockErrorToday: nil,
			mockResultTmrw: price.DailyPriceInfo{
				DayRating:        price.Normal,
				DayAverage:       0.15,
				ThirtyDayAverage: 0.25,
				Prices:           pricesTomorrow,
				ExpensivePeriods: [][]price.Price{
					pricesTomorrow[2:5], pricesTomorrow[8:9],
				},
				CheapPeriods: [][]price.Price{
					pricesTomorrow[4:5],
				},
			},
			mockErrorTmrw:  nil,
			lang:           language.English,
			shouldContain1: "Today is a good day with an average price of 10 cents per kilowatt-hour",
			shouldContain2: "The current price is 24 cents per kilowatt-hour",
			shouldContain3: "The cheap periods for today have already passed",
			shouldContain4: "The expensive periods for today have already passed.",
			shouldContain5: "Tomorrow is a normal day with an average price of 15 cents per kilowatt-hour",
		},
		{
			name: "Full feed - no prices tomorrow",
			t:    time.Date(2023, 1, 1, 23, 0, 0, 0, madridLocation),
			mockResultToday: price.DailyPriceInfo{
				DayRating:        price.Good,
				DayAverage:       0.1,
				ThirtyDayAverage: 0.2,
				Prices:           pricesToday,
				ExpensivePeriods: [][]price.Price{
					pricesToday[2:5], pricesToday[8:9],
				},
				CheapPeriods: [][]price.Price{
					pricesToday[10:11], pricesToday[13:14],
				},
			},
			mockErrorToday: nil,
			mockResultTmrw: price.DailyPriceInfo{},
			mockErrorTmrw:  nil,
			lang:           language.English,
			shouldContain1: "Today is a good day with an average price of 10 cents per kilowatt-hour",
			shouldContain2: "The current price is 24 cents per kilowatt-hour",
			shouldContain3: "The cheap periods for today have already passed",
			shouldContain4: "The expensive periods for today have already passed.",
			shouldContain5: "",
		},
		{
			name: "Full feed - no prices tomorrow error",
			t:    time.Date(2023, 1, 1, 23, 0, 0, 0, madridLocation),
			mockResultToday: price.DailyPriceInfo{
				DayRating:        price.Good,
				DayAverage:       0.1,
				ThirtyDayAverage: 0.2,
				Prices:           pricesToday,
				ExpensivePeriods: [][]price.Price{
					pricesToday[2:5], pricesToday[8:9],
				},
				CheapPeriods: [][]price.Price{
					pricesToday[10:11], pricesToday[13:14],
				},
			},
			mockErrorToday: nil,
			mockResultTmrw: price.DailyPriceInfo{},
			mockErrorTmrw:  errors.New("error"),
			lang:           language.English,
			shouldContain1: "Today is a good day with an average price of 10 cents per kilowatt-hour",
			shouldContain2: "The current price is 24 cents per kilowatt-hour",
			shouldContain3: "The cheap periods for today have already passed",
			shouldContain4: "The expensive periods for today have already passed.",
			shouldContain5: "",
		},
		{
			name:            "Full feed - no prices today",
			t:               time.Date(2023, 1, 1, 23, 0, 0, 0, madridLocation),
			mockResultToday: price.DailyPriceInfo{},
			mockErrorToday:  nil,
			mockResultTmrw:  price.DailyPriceInfo{},
			mockErrorTmrw:   nil,
			lang:            language.English,
			shouldContain1:  "There is no data available yet for today",
			shouldContain2:  "",
			shouldContain3:  "",
			shouldContain4:  "",
			shouldContain5:  "",
		},
		{
			name:            "Full feed - no prices today error",
			t:               time.Date(2023, 1, 1, 23, 0, 0, 0, madridLocation),
			mockResultToday: price.DailyPriceInfo{},
			mockErrorToday:  errors.New("error"),
			mockResultTmrw:  price.DailyPriceInfo{},
			mockErrorTmrw:   nil,
			lang:            language.English,
			shouldContain1:  "",
			shouldContain2:  "",
			shouldContain3:  "",
			shouldContain4:  "",
			shouldContain5:  "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			mockPriceService := &price.MockPriceService{
				MockGetDailyInfoResult: &[]price.DailyPriceInfo{tc.mockResultToday, tc.mockResultTmrw},
				MockGetDailyInfoError:  &[]error{tc.mockErrorToday, tc.mockErrorTmrw},
			}

			service := &Service{
				PriceService: mockPriceService,
			}

			actual, _ := service.GetFullFeed(ctx, tc.t, tc.lang)
			if !strings.Contains(actual, tc.shouldContain1) {
				t.Errorf("'%s' should contain: '%s'", actual, tc.shouldContain1)
			}
			if !strings.Contains(actual, tc.shouldContain2) {
				t.Errorf("'%s' should contain: '%s'", actual, tc.shouldContain2)
			}
			if !strings.Contains(actual, tc.shouldContain3) {
				t.Errorf("'%s' should contain: '%s'", actual, tc.shouldContain3)
			}
			if !strings.Contains(actual, tc.shouldContain4) {
				t.Errorf("'%s' should contain: '%s'", actual, tc.shouldContain4)
			}
			if !strings.Contains(actual, tc.shouldContain5) {
				t.Errorf("'%s' should contain: '%s'", actual, tc.shouldContain5)
			}
		})
	}
}

func TestProcessAlexaSkillRequest(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name                string
		lang                language.Tag
		t                   time.Time
		intent              AlexaIntent
		mockFullFeed        price.DailyPriceInfo
		mockFullFeedError   error
		mockDayRating       price.DayRating
		mockDayRatingError  error
		mockDayAverage      float64
		mockDayAverageError error
		mockThirtyDayAvg    float64
		mockThirtyDayErr    error
		mockGetCheapPeriods [][]price.Price
		mockGetCheapError   error
		mockGetExpensive    [][]price.Price
		mockGetExpensiveErr error
		mockGetPrice        price.Price
		mockGetPriceErr     error
		expectMessage       string
		expectEnd           bool
	}{
		{
			name: "AMAZON.CancelIntent (English)",
			lang: language.English,
			t:    time.Date(2023, 1, 1, 1, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "AMAZON.CancelIntent",
			},
			expectMessage: "Goodbye!",
			expectEnd:     true,
		},
		{
			name: "AMAZON.CancelIntent (Spanish)",
			lang: language.Spanish,
			t:    time.Date(2023, 1, 1, 1, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "AMAZON.CancelIntent",
			},
			expectMessage: "¡Adiós!",
			expectEnd:     true,
		},
		{
			name: "AMAZON.HelpIntent (English)",
			lang: language.English,
			t:    time.Date(2023, 1, 1, 1, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "AMAZON.HelpIntent",
			},
			expectMessage: "You can ask for the current price, the average price today, the average price for the last 30 days, the next cheap period, the next expensive period, a full update or the prices for tomorrow. What would you like to know?",
			expectEnd:     false,
		},
		{
			name: "AMAZON.HelpIntent (Spanish)",
			lang: language.Spanish,
			t:    time.Date(2023, 1, 1, 1, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "AMAZON.HelpIntent",
			},
			expectMessage: "Esta skill le permite obtener información sobre el precio de la electricidad en España. Puede preguntar por el precio actual, el precio promedio de hoy, el precio promedio de los últimos 30 días, el próximo período barato, el próximo período caro, una actualización completa o por los precio de mañana. Que le gustaría saber?",
			expectEnd:     false,
		},
		{
			name: "AMAZON.StopIntent (English)",
			lang: language.English,
			t:    time.Date(2023, 1, 1, 1, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "AMAZON.StopIntent",
			},
			expectMessage: "Goodbye!",
			expectEnd:     true,
		},
		{
			name: "AMAZON.StopIntent (Spanish)",
			lang: language.Spanish,
			t:    time.Date(2023, 1, 1, 1, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "AMAZON.StopIntent",
			},
			expectMessage: "¡Adiós!",
			expectEnd:     true,
		},
		{
			name: "AMAZON.NavigateHomeIntent (English)",
			lang: language.English,
			t:    time.Date(2023, 1, 1, 1, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "AMAZON.NavigateHomeIntent",
			},
			expectMessage: "Welcome to the electricity prices skill. Say, give me a full update.",
			expectEnd:     false,
		},
		{
			name: "AMAZON.NavigateHomeIntent (Spanish)",
			lang: language.Spanish,
			t:    time.Date(2023, 1, 1, 1, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "AMAZON.NavigateHomeIntent",
			},
			expectMessage: "Bienvenido a la skill de precios de la electricidad. Diga, dame una actualización completa.",
			expectEnd:     false,
		},
		{
			name: "AMAZON.FallbackIntent (English)",
			lang: language.English,
			t:    time.Date(2023, 1, 1, 1, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "AMAZON.FallbackIntent",
			},
			expectMessage: "Welcome to the electricity prices skill. Say, give me a full update.",
			expectEnd:     false,
		},
		{
			name: "AMAZON.FallbackIntent (Spanish)",
			lang: language.Spanish,
			t:    time.Date(2023, 1, 1, 1, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "AMAZON.FallbackIntent",
			},
			expectMessage: "Bienvenido a la skill de precios de la electricidad. Diga, dame una actualización completa.",
			expectEnd:     false,
		},
		{
			name: "FULL (English)",
			lang: language.English,
			t:    time.Date(2023, 1, 1, 1, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "FULL",
			},
			mockFullFeed: price.DailyPriceInfo{
				DayRating:        price.Good,
				DayAverage:       0.1,
				ThirtyDayAverage: 0.2,
				Prices:           pricesToday,
				ExpensivePeriods: [][]price.Price{
					pricesToday[2:5], pricesToday[8:9],
				},
				CheapPeriods: [][]price.Price{
					pricesToday[10:11], pricesToday[13:14],
				},
			},
			mockFullFeedError: nil,
			expectMessage:     "Today is a good day with an average price of 10 cents per kilowatt-hour. The current price is 2 cents per kilowatt-hour. The next cheap period starts at 10 AM with an average price of 11 cents per kilowatt-hour and will end at 11 AM. The next expensive period starts at 2 AM with an average price of 4 cents per kilowatt-hour and will end at 5 AM.",
			expectEnd:         false,
		},
		{
			name: "FULL (Spanish)",
			lang: language.Spanish,
			t:    time.Date(2023, 1, 1, 1, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "FULL",
			},
			mockFullFeed: price.DailyPriceInfo{
				DayRating:        price.Good,
				DayAverage:       0.1,
				ThirtyDayAverage: 0.2,
				Prices:           pricesToday,
				ExpensivePeriods: [][]price.Price{
					pricesToday[2:5], pricesToday[8:9],
				},
				CheapPeriods: [][]price.Price{
					pricesToday[10:11], pricesToday[13:14],
				},
			},
			mockFullFeedError: nil,
			expectMessage:     "Hoy es un día bueno, con un precio medio de 10 céntimos por kilovatio-hora. El precio actual es 2 céntimos por kilovatio-hora. El próximo período barato comienza a las 10 AM con un precio medio de 11 céntimos por kilovatio-hora y terminara a las 11 AM. El próximo período caro comienza a las 2 AM con un precio promedio de 4 céntimos por kilovatio-hora y terminara a las 5 AM.",
			expectEnd:         false,
		},
		{
			name: "FULL - error",
			lang: language.English,
			t:    time.Date(2023, 1, 1, 1, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "FULL",
			},
			mockFullFeed:      price.DailyPriceInfo{},
			mockFullFeedError: errors.New("error"),
			expectMessage:     "Sorry, there was an error. Please try again later.",
			expectEnd:         false,
		},
		{
			name: "TODAY (English)",
			lang: language.English,
			t:    time.Date(2023, 1, 1, 3, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "TODAY",
			},
			mockDayRating:       price.Good,
			mockDayRatingError:  nil,
			mockDayAverage:      0.1,
			mockDayAverageError: nil,
			expectMessage:       "Today is a good day with an average price of 10 cents per kilowatt-hour.",
			expectEnd:           false,
		},
		{
			name: "TODAY (Spanish)",
			lang: language.Spanish,
			t:    time.Date(2023, 1, 1, 3, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "TODAY",
			},
			mockDayRating:       price.Bad,
			mockDayRatingError:  nil,
			mockDayAverage:      0.2,
			mockDayAverageError: nil,
			expectMessage:       "Hoy es un día malo, con un precio medio de 20 céntimos por kilovatio-hora.",
			expectEnd:           false,
		},
		{
			name: "TODAY - error (English)",
			lang: language.English,
			t:    time.Date(2023, 1, 1, 3, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "TODAY",
			},
			mockDayRating:       price.Good,
			mockDayRatingError:  nil,
			mockDayAverage:      0.1,
			mockDayAverageError: errors.New("error"),
			expectMessage:       "Sorry, there was an error. Please try again later.",
			expectEnd:           false,
		},
		{
			name: "TODAY - error (Spanish)",
			lang: language.Spanish,
			t:    time.Date(2023, 1, 1, 3, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "TODAY",
			},
			mockDayRating:       price.Normal,
			mockDayRatingError:  nil,
			mockDayAverage:      0.1,
			mockDayAverageError: errors.New("error"),
			expectMessage:       "Lo siento, no pude obtener los datos. Por favor, inténtelo de nuevo más tarde.",
			expectEnd:           false,
		},
		{
			name: "TODAY_AVERAGE (English)",
			lang: language.English,
			t:    time.Date(2023, 1, 1, 3, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "TODAY_AVERAGE",
			},
			mockDayRating:       price.Good,
			mockDayRatingError:  nil,
			mockDayAverage:      0.1,
			mockDayAverageError: nil,
			expectMessage:       "Today is a good day with an average price of 10 cents per kilowatt-hour.",
			expectEnd:           false,
		},
		{
			name: "TODAY_AVERAGE (Spanish)",
			lang: language.Spanish,
			t:    time.Date(2023, 1, 1, 3, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "TODAY_AVERAGE",
			},
			mockDayRating:       price.Bad,
			mockDayRatingError:  nil,
			mockDayAverage:      0.2,
			mockDayAverageError: nil,
			expectMessage:       "Hoy es un día malo, con un precio medio de 20 céntimos por kilovatio-hora.",
			expectEnd:           false,
		},
		{
			name: "TODAY_AVERAGE - error (English)",
			lang: language.English,
			t:    time.Date(2023, 1, 1, 3, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "TODAY_AVERAGE",
			},
			mockDayRating:       price.Good,
			mockDayRatingError:  nil,
			mockDayAverage:      0.1,
			mockDayAverageError: errors.New("error"),
			expectMessage:       "Sorry, there was an error. Please try again later.",
			expectEnd:           false,
		},
		{
			name: "TODAY_AVERAGE - error (Spanish)",
			lang: language.Spanish,
			t:    time.Date(2023, 1, 1, 3, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "TODAY_AVERAGE",
			},
			mockDayRating:       price.Normal,
			mockDayRatingError:  nil,
			mockDayAverage:      0.1,
			mockDayAverageError: errors.New("error"),
			expectMessage:       "Lo siento, no pude obtener los datos. Por favor, inténtelo de nuevo más tarde.",
			expectEnd:           false,
		},
		{
			name: "TOMORROW (English)",
			lang: language.English,
			t:    time.Date(2023, 1, 1, 3, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "TOMORROW",
			},
			mockDayRating:       price.Normal,
			mockDayRatingError:  nil,
			mockDayAverage:      0.1,
			mockDayAverageError: nil,
			expectMessage:       "Tomorrow is a normal day with an average price of 10 cents per kilowatt-hour.",
			expectEnd:           false,
		},
		{
			name: "TOMORROW (Spanish)",
			lang: language.Spanish,
			t:    time.Date(2023, 1, 1, 3, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "TOMORROW",
			},
			mockDayRating:       price.Bad,
			mockDayRatingError:  nil,
			mockDayAverage:      0.2,
			mockDayAverageError: nil,
			expectMessage:       "Mañana es un día malo, con un precio promedio de 20 céntimos por kilovatio-hora.",
			expectEnd:           false,
		},
		{
			name: "TOMORROW - error (English)",
			lang: language.English,
			t:    time.Date(2023, 1, 1, 3, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "TOMORROW",
			},
			mockDayRating:       price.Good,
			mockDayRatingError:  nil,
			mockDayAverage:      0.1,
			mockDayAverageError: errors.New("error"),
			expectMessage:       "There is no data available yet for tomorrow. Please check back later. Prices are generally available by 8:30 PM.",
			expectEnd:           false,
		},
		{
			name: "TOMORROW - error (Spanish)",
			lang: language.Spanish,
			t:    time.Date(2023, 1, 1, 3, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "TOMORROW",
			},
			mockDayRating:       price.Normal,
			mockDayRatingError:  nil,
			mockDayAverage:      0.1,
			mockDayAverageError: errors.New("error"),
			expectMessage:       "Aún no hay datos disponibles para mañana. Por favor, vuelva más tarde. Los precios están generalmente disponibles a las 8:30 PM.",
			expectEnd:           false,
		},
		{
			name: "NEXT_CHEAP (English)",
			lang: language.English,
			t:    time.Date(2023, 1, 1, 10, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "NEXT_CHEAP",
			},
			mockGetCheapPeriods: [][]price.Price{
				pricesToday[10:11], pricesToday[13:14],
			},
			mockGetCheapError: nil,
			expectMessage:     "You are currently in a cheap period that started at 10 AM with an average price of 11 cents per kilowatt-hour and will end at 11 AM.",
			expectEnd:         false,
		},
		{
			name: "NEXT_CHEAP (Spanish)",
			lang: language.Spanish,
			t:    time.Date(2023, 1, 1, 10, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "NEXT_CHEAP",
			},
			mockGetCheapPeriods: [][]price.Price{
				pricesToday[10:11], pricesToday[12:15],
			},
			mockGetCheapError: nil,
			expectMessage:     "Actualmente se encuentra en un período barato que comenzó a las 10 AM con un precio promedio de 11 céntimos por kilovatio-hora y terminara a las 11 AM.",
			expectEnd:         false,
		},
		{
			name: "NEXT_CHEAP - error (English)",
			lang: language.English,
			t:    time.Date(2023, 1, 1, 10, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "NEXT_CHEAP",
			},
			mockGetCheapPeriods: [][]price.Price{
				pricesToday[10:11], pricesToday[13:14],
			},
			mockGetCheapError: errors.New("error"),
			expectMessage:     "Sorry, there was an error. Please try again later.",
			expectEnd:         false,
		},
		{
			name: "NEXT_CHEAP - error (Spanish)",
			lang: language.Spanish,
			t:    time.Date(2023, 1, 1, 10, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "NEXT_CHEAP",
			},
			mockGetCheapPeriods: nil,
			mockGetCheapError:   errors.New("error"),
			expectMessage:       "Lo siento, no pude obtener los datos. Por favor, inténtelo de nuevo más tarde.",
			expectEnd:           false,
		},
		{
			name: "NEXT_EXPENSIVE (English)",
			lang: language.English,
			t:    time.Date(2023, 1, 1, 10, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "NEXT_EXPENSIVE",
			},
			mockGetExpensive: [][]price.Price{
				pricesToday[2:5], pricesToday[8:9],
			},
			mockGetExpensiveErr: nil,
			expectMessage:       "The expensive periods for today have already passed.",
			expectEnd:           false,
		},
		{
			name: "NEXT_EXPENSIVE (Spanish)",
			lang: language.Spanish,
			t:    time.Date(2023, 1, 1, 10, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "NEXT_EXPENSIVE",
			},
			mockGetExpensive: [][]price.Price{
				pricesToday[2:5], pricesToday[8:9],
			},
			mockGetExpensiveErr: nil,
			expectMessage:       "Los períodos caros de hoy ya han pasado.",
			expectEnd:           false,
		},
		{
			name: "NEXT_EXPENSIVE - error (English)",
			lang: language.English,
			t:    time.Date(2023, 1, 1, 10, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "NEXT_EXPENSIVE",
			},
			mockGetExpensive: [][]price.Price{
				pricesToday[2:5], pricesToday[8:9],
			},
			mockGetExpensiveErr: errors.New("error"),
			expectMessage:       "Sorry, there was an error. Please try again later.",
			expectEnd:           false,
		},
		{
			name: "NEXT_EXPENSIVE - error (Spanish)",
			lang: language.Spanish,
			t:    time.Date(2023, 1, 1, 10, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "NEXT_EXPENSIVE",
			},
			mockGetExpensive: [][]price.Price{
				pricesToday[2:5], pricesToday[8:9],
			},
			mockGetExpensiveErr: errors.New("error"),
			expectMessage:       "Lo siento, no pude obtener los datos. Por favor, inténtelo de nuevo más tarde.",
			expectEnd:           false,
		},
		{
			name: "CURRENT_PRICE (English)",
			lang: language.English,
			t:    time.Date(2023, 1, 1, 10, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "CURRENT_PRICE",
			},
			mockGetPrice:    pricesToday[10],
			mockGetPriceErr: nil,
			expectMessage:   "The current price is 11 cents per kilowatt-hour.",
			expectEnd:       false,
		},
		{
			name: "CURRENT_PRICE (Spanish)",
			lang: language.Spanish,
			t:    time.Date(2023, 1, 1, 10, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "CURRENT_PRICE",
			},
			mockGetPrice:    pricesToday[9],
			mockGetPriceErr: nil,
			expectMessage:   "El precio actual es 10 céntimos por kilovatio-hora.",
			expectEnd:       false,
		},
		{
			name: "CURRENT_PRICE - error (English)",
			lang: language.English,
			t:    time.Date(2023, 1, 1, 10, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "CURRENT_PRICE",
			},
			mockGetPrice:    pricesToday[10],
			mockGetPriceErr: errors.New("error"),
			expectMessage:   "Sorry, there was an error. Please try again later.",
			expectEnd:       false,
		},
		{
			name: "CURRENT_PRICE - error (Spanish)",
			lang: language.Spanish,
			t:    time.Date(2023, 1, 1, 10, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "CURRENT_PRICE",
			},
			mockGetPrice:    pricesToday[9],
			mockGetPriceErr: errors.New("error"),
			expectMessage:   "Lo siento, no pude obtener los datos. Por favor, inténtelo de nuevo más tarde.",
			expectEnd:       false,
		},
		{
			name: "THIRTY_DAY_AVERAGE (English)",
			lang: language.English,
			t:    time.Date(2023, 1, 1, 10, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "THIRTY_DAY_AVERAGE",
			},
			mockThirtyDayAvg: 0.1,
			mockThirtyDayErr: nil,
			expectMessage:    "The average price for the last 30 days is 10 cents per kilowatt-hour.",
			expectEnd:        false,
		},
		{
			name: "THIRTY_DAY_AVERAGE (Spanish)",
			lang: language.Spanish,
			t:    time.Date(2023, 1, 1, 10, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "THIRTY_DAY_AVERAGE",
			},
			mockThirtyDayAvg: 0.2,
			mockThirtyDayErr: nil,
			expectMessage:    "El precio promedio de los últimos 30 días es 20 céntimos por kilovatio-hora.",
			expectEnd:        false,
		},
		{
			name: "THIRTY_DAY_AVERAGE - error (English)",
			lang: language.English,
			t:    time.Date(2023, 1, 1, 10, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "THIRTY_DAY_AVERAGE",
			},
			mockThirtyDayAvg: 0.1,
			mockThirtyDayErr: errors.New("error"),
			expectMessage:    "Sorry, there was an error. Please try again later.",
			expectEnd:        false,
		},
		{
			name: "THIRTY_DAY_AVERAGE - error (Spanish)",
			lang: language.Spanish,
			t:    time.Date(2023, 1, 1, 10, 0, 0, 0, madridLocation),
			intent: AlexaIntent{
				Name: "THIRTY_DAY_AVERAGE",
			},
			mockThirtyDayAvg: 0.2,
			mockThirtyDayErr: errors.New("error"),
			expectMessage:    "Lo siento, no pude obtener los datos. Por favor, inténtelo de nuevo más tarde.",
			expectEnd:        false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockPriceService := &price.MockPriceService{
				MockGetDailyInfoResult:        &[]price.DailyPriceInfo{tc.mockFullFeed},
				MockGetDailyInfoError:         &[]error{tc.mockFullFeedError},
				MockGetDayRatingResult:        &[]price.DayRating{tc.mockDayRating},
				MockGetDayRatingError:         &[]error{tc.mockDayRatingError},
				MockGetDayAverageResult:       &[]float64{tc.mockDayAverage},
				MockGetDayAverageError:        &[]error{tc.mockDayAverageError},
				MockGetCheapPeriodsResult:     &[][][]price.Price{tc.mockGetCheapPeriods},
				MockGetCheapPeriodsError:      &[]error{tc.mockGetCheapError},
				MockGetExpensivePeriodsResult: &[][][]price.Price{tc.mockGetExpensive},
				MockGetExpensivePeriodsError:  &[]error{tc.mockGetExpensiveErr},
				MockGetPriceResult:            &[]price.Price{tc.mockGetPrice},
				MockGetPriceError:             &[]error{tc.mockGetPriceErr},
				MockGetThirtyDayAverageResult: &[]float64{tc.mockThirtyDayAvg},
				MockGetThirtyDayAverageError:  &[]error{tc.mockThirtyDayErr},
			}

			service := &Service{
				PriceService: mockPriceService,
			}

			res := service.ProcessAlexaSkillRequest(ctx, tc.intent, tc.t, tc.lang)
			if res.Response.OutputSpeech.Text != tc.expectMessage {
				t.Errorf("expected '%s' but got '%s'", tc.expectMessage, res.Response.OutputSpeech.Text)
			}
			if res.Response.ShouldEndSession != tc.expectEnd {
				t.Errorf("expected '%v' but got '%v'", tc.expectEnd, res.Response.ShouldEndSession)
			}
		})
	}
}
