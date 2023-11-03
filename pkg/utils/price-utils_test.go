package utils

import (
	"electricity-prices/pkg/model"
	"encoding/json"
	"math"
	"os"
	"testing"
	"time"
)

const epsilon = 1e-4 // Tolerance level

func floatEquals(a, b float64) bool {
	return math.Abs(a-b) <= epsilon
}

func readJSONFromFile(filePath string, result interface{}) error {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, result)
	if err != nil {
		return err
	}
	return nil
}

func TestCalculateAverage(t *testing.T) {
	testCases := []struct {
		name     string
		prices   []model.Price
		expected float64
	}{
		{"Empty slice", []model.Price{}, 0.0},
		{"One price", []model.Price{{Price: 1.0}}, 1.0},
		{"Two prices", []model.Price{{Price: 1.0}, {Price: 2.0}}, 1.5},
		{"Three prices", []model.Price{{Price: 1.0}, {Price: 2.0}, {Price: 3.0}}, 2.0},
		{"Mixed order", []model.Price{{Price: 3.0}, {Price: 1.0}, {Price: 2.0}}, 2.0},
		{"Negative", []model.Price{{Price: -1.0}, {Price: 2.0}, {Price: 3.0}}, 1.333333},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			average := CalculateAverage(tc.prices)
			if !floatEquals(average, tc.expected) {
				t.Errorf("Expected %f, but got %f", tc.expected, average)
			}
		})
	}
}

func TestSortDailyAverages(t *testing.T) {
	averages := []model.DailyAverage{
		{Date: "2020-01-01"},
		{Date: "2020-01-03"},
		{Date: "2020-01-02"},
	}

	sortDailyAverages(averages)

	if averages[0].Date != "2020-01-01" {
		t.Errorf("Expected first average to be 2020-01-01 but was %s", averages[0].Date)
	}

	if averages[1].Date != "2020-01-02" {
		t.Errorf("Expected second average to be 2020-01-02 but was %s", averages[1].Date)
	}

	if averages[2].Date != "2020-01-03" {
		t.Errorf("Expected third average to be 2020-01-03 but was %s", averages[2].Date)
	}
}

func TestCalculateDayRating(t *testing.T) {
	testCases := []struct {
		name         string
		dayAvg       float64
		thirtyDayAvg float64
		expected     model.DayRating
	}{
		{"Good", 0.2, 0.1, model.Bad},
		{"Bad", 0.1, 0.2, model.Good},
		{"Normal", 0.1, 0.11, model.Normal},
		{"Zeros", 0.0, 0.0, model.Normal},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rating := CalculateDayRating(tc.dayAvg, tc.thirtyDayAvg)
			if rating != tc.expected {
				t.Errorf("Expected %s, but got %s", tc.expected, rating)
			}
		})
	}
}

func TestCalculateCombinedAverage(t *testing.T) {
	testCases := []struct {
		name         string
		dayAvg       float64
		thirtyDayAvg float64
		expected     float64
	}{
		{"Both zero", 0.0, 0.0, 0.0},
		{"Day zero", 0.0, 1.0, 0.333333},
		{"Thirty day zero", 1.0, 0.0, 0.666667},
		{"Both non-zero", 1.0, 1.0, 1.0},
		{"Both negative", -1.0, -1.0, -1.0},
		{"Day negative", -1.0, 1.0, -0.333333},
		{"Thirty day negative", 1.0, -1.0, 0.333333},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			average := calculateCombinedAverage(tc.dayAvg, tc.thirtyDayAvg)
			if !floatEquals(average, tc.expected) {
				t.Errorf("Expected %f, but got %f", tc.expected, average)
			}
		})
	}
}

func TestGetMinPrice(t *testing.T) {
	testCases := []struct {
		name     string
		prices   []model.Price
		expected float64
	}{
		{"Empty slice", []model.Price{}, 0.0},
		{"One price", []model.Price{{Price: 1.0}}, 1.0},
		{"Two prices", []model.Price{{Price: 1.0}, {Price: 2.0}}, 1.0},
		{"Three prices", []model.Price{{Price: 1.0}, {Price: 2.0}, {Price: 3.0}}, 1.0},
		{"Mixed order", []model.Price{{Price: 3.0}, {Price: 1.0}, {Price: 2.0}}, 1.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			minP := getMinPrice(tc.prices)
			if !floatEquals(minP, tc.expected) {
				t.Errorf("Expected %f, but got %f", tc.expected, minP)
			}
		})
	}
}

func TestGetMinAndMaxPrices(t *testing.T) {
	testCases := []struct {
		name        string
		prices      []model.Price
		expectedMin float64
		expectedMax float64
	}{
		{"Empty slice", []model.Price{}, 0.0, 0.0},
		{"One price", []model.Price{{Price: 1.0}}, 1.0, 1.0},
		{"Two prices", []model.Price{{Price: 1.0}, {Price: 2.0}}, 1.0, 2.0},
		{"Three prices", []model.Price{{Price: 1.0}, {Price: 2.0}, {Price: 3.0}}, 1.0, 3.0},
		{"Mixed order", []model.Price{{Price: 3.0}, {Price: 1.0}, {Price: 2.0}}, 1.0, 3.0},
		{"Negative", []model.Price{{Price: -1.0}, {Price: 2.0}, {Price: 3.0}}, -1.0, 3.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			minP, maxP := getMinAndMaxPrices(tc.prices)
			if !floatEquals(minP, tc.expectedMin) {
				t.Errorf("Expected %f, but got %f", tc.expectedMin, minP)
			}
			if !floatEquals(maxP, tc.expectedMax) {
				t.Errorf("Expected %f, but got %f", tc.expectedMax, maxP)
			}
		})
	}
}

func TestGetMaxPrice(t *testing.T) {
	testCases := []struct {
		name     string
		prices   []model.Price
		expected float64
	}{
		{"Empty slice", []model.Price{}, 0.0},
		{"One price", []model.Price{{Price: 1.0}}, 1.0},
		{"Two prices", []model.Price{{Price: 1.0}, {Price: 2.0}}, 2.0},
		{"Three prices", []model.Price{{Price: 1.0}, {Price: 2.0}, {Price: 3.0}}, 3.0},
		{"Mixed order", []model.Price{{Price: 3.0}, {Price: 1.0}, {Price: 2.0}}, 3.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			maxP := getMaxPrice(tc.prices)
			if !floatEquals(maxP, tc.expected) {
				t.Errorf("Expected %f, but got %f", tc.expected, maxP)
			}
		})
	}
}

func TestCalculateMinVariance(t *testing.T) {
	testCases := []struct {
		name     string
		minPrice float64
		maxPrice float64
		expected float64
	}{
		{"Zeroes", 0.0, 0.0, 0.0},
		{"One", 1.0, 1.0, 0.0},
		{"Two", 1.0, 2.0, 0.166667},
		{"Three", 1.0, 3.0, 0.333333},
		{"Negative", -1.0, 1.0, 0.333333},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			variance := calculateMinVariance(tc.minPrice, tc.maxPrice)
			if !floatEquals(variance, tc.expected) {
				t.Errorf("Expected %f, but got %f", tc.expected, variance)
			}
		})
	}
}

func TestCalculateMaxVariance(t *testing.T) {
	testCases := []struct {
		name     string
		minPrice float64
		maxPrice float64
		expected float64
	}{
		{"Zeroes", 0.0, 0.0, 0.0},
		{"One", 1.0, 1.0, 0.0},
		{"Two", 1.0, 2.0, 0.333333},
		{"Three", 1.0, 3.0, 0.666667},
		{"Negative", -1.0, 1.0, 0.666667},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			variance := calculateMaxVariance(tc.minPrice, tc.maxPrice)
			if !floatEquals(variance, tc.expected) {
				t.Errorf("Expected %f, but got %f", tc.expected, variance)
			}
		})
	}
}

func TestCalculateCheapVariance(t *testing.T) {
	var normalPeriod []model.Price
	err := readJSONFromFile("../../test/resources/normal-day.json", &normalPeriod)
	if err != nil {
		t.Errorf("Error reading file: %s", err)
	}

	testCases := []struct {
		name         string
		prices       []model.Price
		thirtyDayAvg float64
		expected     float64
	}{
		{"Empty slice", []model.Price{}, 1.0, 0.0},
		{"One price", []model.Price{{Price: 1.0}}, 1.0, 0.0},
		{"Two prices", []model.Price{{Price: 1.0}, {Price: 2.0}}, 1.0, 0.166667},
		{"Three prices", []model.Price{{Price: 1.0}, {Price: 2.0}, {Price: 3.0}}, 1.0, 0.333333},
		{"Negative", []model.Price{{Price: -1.0}, {Price: 2.0}, {Price: 3.0}}, 1.0, 1.111111},
		{"normal day", normalPeriod, 0.15, 0.016479999999999998},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			variance := calculateCheapVariance(tc.prices, tc.thirtyDayAvg)
			if !floatEquals(variance, tc.expected) {
				t.Errorf("Expected %f, but got %f", tc.expected, variance)
			}
		})
	}
}

func TestCalculateExpensiveVariance(t *testing.T) {
	var normalPeriod []model.Price
	err := readJSONFromFile("../../test/resources/normal-day.json", &normalPeriod)
	if err != nil {
		t.Errorf("Error reading file: %s", err)
	}

	testCases := []struct {
		name         string
		prices       []model.Price
		thirtyDayAvg float64
		expected     float64
	}{
		{"Empty slice", []model.Price{}, 1.0, 0.0},
		{"One price", []model.Price{{Price: 1.0}}, 1.0, 0.0},
		{"Two prices", []model.Price{{Price: 1.0}, {Price: 2.0}}, 1.0, 0.333333},
		{"Three prices", []model.Price{{Price: 1.0}, {Price: 2.0}, {Price: 3.0}}, 1.0, 0.666667},
		{"Negative", []model.Price{{Price: -1.0}, {Price: 2.0}, {Price: 3.0}}, 1.0, 0.888889},
		{"normal day", normalPeriod, 0.15, 0.032959999999999996},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			variance := calculateExpensiveVariance(tc.prices, tc.thirtyDayAvg)
			if !floatEquals(variance, tc.expected) {
				t.Errorf("Expected %f, but got %f", tc.expected, variance)
			}
		})
	}
}

func TestCalculateCheapPeriods(t *testing.T) {
	var cheapPeriod []model.Price
	err := readJSONFromFile("../../test/resources/cheap-day.json", &cheapPeriod)
	if err != nil {
		t.Errorf("Error reading file: %s", err)
	}

	var cheapPeriodResult [][]model.Price
	err = readJSONFromFile("../../test/resources/cheap-day-cp.json", &cheapPeriodResult)
	if err != nil {
		t.Errorf("Error reading file: %s", err)
	}

	testCases := []struct {
		name         string
		prices       []model.Price
		thirtyDayAvg float64
		expected     [][]model.Price
	}{
		{"Empty slice", []model.Price{}, 1.0, [][]model.Price{}},
		{"One price", []model.Price{{Price: 1.0}}, 1.0, [][]model.Price{{model.Price{Price: 1.0}}}},
		{"Two prices", []model.Price{{Price: 1.0}, {Price: 2.0}}, 1.0, [][]model.Price{{model.Price{Price: 1.0}}}},
		{"Three prices", []model.Price{{Price: 1.0}, {Price: 2.0}, {Price: 3.0}}, 1.0, [][]model.Price{{model.Price{Price: 1.0}}}},
		{"Negative", []model.Price{{Price: -1.0}, {Price: 2.0}, {Price: 3.0}}, 1.0, [][]model.Price{{model.Price{Price: -1.0}}}},
		{"periodOne", cheapPeriod, 0.15, cheapPeriodResult},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cheapPeriods := CalculateCheapPeriods(tc.prices, tc.thirtyDayAvg)
			if len(cheapPeriods) != len(tc.expected) {
				t.Errorf("Expected %d periods, but got %d", len(tc.expected), len(cheapPeriods))
			}
			for i := range cheapPeriods {
				if len(cheapPeriods[i]) != len(tc.expected[i]) {
					t.Errorf("Expected %d prices in period %d, but got %d", len(tc.expected[i]), i, len(cheapPeriods[i]))
				}
				for j := range cheapPeriods[i] {
					if !floatEquals(cheapPeriods[i][j].Price, tc.expected[i][j].Price) {
						t.Errorf("Expected %f, but got %f", tc.expected[i][j].Price, cheapPeriods[i][j].Price)
					}
				}
			}
		})
	}
}

func TestCalculateExpensivePeriods(t *testing.T) {
	var cheapPeriod []model.Price
	err := readJSONFromFile("../../test/resources/cheap-day.json", &cheapPeriod)
	if err != nil {
		t.Errorf("Error reading file: %s", err)
	}

	testCases := []struct {
		name         string
		prices       []model.Price
		thirtyDayAvg float64
		expected     [][]model.Price
	}{
		{"Empty slice", []model.Price{}, 1.0, [][]model.Price{}},
		{"One price", []model.Price{{Price: 1.0}}, 1.0, [][]model.Price{{model.Price{Price: 1.0}}}},
		{"Two prices", []model.Price{{Price: 1.0}, {Price: 2.0}}, 1.0, [][]model.Price{{model.Price{Price: 2.0}}}},
		{"Three prices", []model.Price{{Price: 1.0}, {Price: 2.0}, {Price: 3.0}}, 1.0, [][]model.Price{{model.Price{Price: 3.0}}}},
		{"Negative", []model.Price{{Price: -1.0}, {Price: 2.0}, {Price: 3.0}}, 1.0, [][]model.Price{{model.Price{Price: 3.0}}}},
		{"periodOne", cheapPeriod, 0.15, [][]model.Price{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expensivePeriods := CalculateExpensivePeriods(tc.prices, tc.thirtyDayAvg)
			if len(expensivePeriods) != len(tc.expected) {
				t.Errorf("Expected %d periods, but got %d", len(tc.expected), len(expensivePeriods))
			}
			for i := range expensivePeriods {
				if len(expensivePeriods[i]) != len(tc.expected[i]) {
					t.Errorf("Expected %d prices in period %d, but got %d", len(tc.expected[i]), i, len(expensivePeriods[i]))
				}
				for j := range expensivePeriods[i] {
					if !floatEquals(expensivePeriods[i][j].Price, tc.expected[i][j].Price) {
						t.Errorf("Expected %f, but got %f", tc.expected[i][j].Price, expensivePeriods[i][j].Price)
					}
				}
			}
		})
	}
}

func TestGetNextPeriod(t *testing.T) {

	date := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	price1 := model.Price{Price: 1.0, DateTime: date}
	price2 := model.Price{Price: 2.0, DateTime: date.Add(time.Hour)}
	price3 := model.Price{Price: 3.0, DateTime: date.Add(2 * time.Hour)}
	price4 := model.Price{Price: 4.0, DateTime: date.Add(3 * time.Hour)}
	price5 := model.Price{Price: 5.0, DateTime: date.Add(4 * time.Hour)}
	price6 := model.Price{Price: 6.0, DateTime: date.Add(5 * time.Hour)}
	price7 := model.Price{Price: 7.0, DateTime: date.Add(6 * time.Hour)}

	testCases := []struct {
		name            string
		periods         [][]model.Price
		date            time.Time
		expectedP       []model.Price
		expectedStarted bool
	}{
		{"Empty slice", [][]model.Price{}, date, []model.Price{}, false},
		{"One period in future", [][]model.Price{{price2}}, date, []model.Price{price2}, false},
		{"One period in past", [][]model.Price{{price2}}, date.Add(5 * time.Hour), nil, false},
		{"One period that has started", [][]model.Price{{price1, price2, price3, price4}}, date.Add(time.Hour), []model.Price{price1, price2, price3, price4}, true},
		{"Two periods both in future", [][]model.Price{{price2}, {price5, price6, price7}}, date, []model.Price{price2}, false},
		{"Two periods one in future", [][]model.Price{{price2}, {price6, price7}}, date.Add(3 * time.Hour), []model.Price{price6, price7}, false},
		{"Three periods both in past", [][]model.Price{{price2}, {price5, price6, price7}}, date.Add(8 * time.Hour), nil, false},
		{"Three periods one in the middle", [][]model.Price{{price2}, {price5, price6, price7}}, date.Add(6 * time.Hour), []model.Price{price5, price6, price7}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			period, started := GetNextPeriod(tc.periods, tc.date)
			if len(period) != len(tc.expectedP) {
				t.Errorf("Expected %d prices, but got %d", len(tc.expectedP), len(period))
			}
			for i := range period {
				if !floatEquals(period[i].Price, tc.expectedP[i].Price) {
					t.Errorf("Expected %f, but got %f", tc.expectedP[i].Price, period[i].Price)
				}
			}
			if started != tc.expectedStarted {
				t.Errorf("Expected %t, but got %t", tc.expectedStarted, started)
			}
		})
	}

}

func TestFormatPrice(t *testing.T) {
	testCases := []struct {
		name     string
		price    float64
		expected string
	}{
		{"Zero", 0.0, "0"},
		{"One", 0.01, "1"},
		{"Two", 0.02, "2"},
		{"Three", 0.03, "3"},
		{"1 decimal place", 0.123, "12.3"},
		{"2 decimal places", 0.1234, "12.3"},
		{"3 decimal places", 0.12345, "12.3"},
		{"4 decimal places", 0.123456, "12.3"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := FormatPrice(tc.price)
			if result != tc.expected {
				t.Errorf("Expected %s, but got %s", tc.expected, result)
			}
		})
	}
}
