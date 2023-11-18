package price

import (
	"context"
	"electricity-prices/pkg/date"
	"fmt"
	"time"
)

func GetDailyPrices(ctx context.Context, t time.Time) ([]Price, error) {
	start, end := date.ParseStartAndEndTimes(t, 1)

	prices, err := getPrices(ctx, start, end)

	if err != nil {
		return nil, err
	}

	return prices, nil
}

func GetDailyAverages(ctx context.Context, date time.Time, numberOfDays int) ([]DailyAverage, error) {

	xDaysAgo := date.AddDate(0, 0, -numberOfDays)
	nextDay := time.Date(date.Year(), date.Month(), date.Day()+1, 0, 0, 0, 0, date.Location())

	// Subtract one second to get the last second of the current day
	today := nextDay.Add(-time.Second)

	prices, err := getPrices(ctx, xDaysAgo, today)

	if err != nil {
		return nil, err
	}

	averages := CalculateDailyAverages(prices)

	return averages, nil

}

func GetDailyInfo(ctx context.Context, date time.Time) (DailyPriceInfo, error) {
	// Get the prices for the given day
	prices, err := GetDailyPrices(ctx, date)
	if err != nil {
		return DailyPriceInfo{}, err
	}

	// Get thirty-day average
	avgPrice, err := GetThirtyDayAverage(ctx, date)
	if err != nil {
		return DailyPriceInfo{}, err
	}

	// Get day rating
	dayAvg := CalculateAverage(prices)
	dayRating := CalculateDayRating(dayAvg, avgPrice)

	// Get cheap periods
	cheapPeriods := CalculateCheapPeriods(prices, avgPrice)

	// Get expensive periods
	expensivePeriods := CalculateExpensivePeriods(prices, avgPrice)

	return DailyPriceInfo{
		Prices:           prices,
		ThirtyDayAverage: avgPrice,
		DayRating:        dayRating,
		DayAverage:       dayAvg,
		CheapPeriods:     cheapPeriods,
		ExpensivePeriods: expensivePeriods,
	}, nil
}

func GetDayRating(ctx context.Context, date time.Time) (DayRating, error) {
	// Get the prices for the given day
	prices, err := GetDailyPrices(ctx, date)
	if err != nil {
		return Nil, err
	}
	if len(prices) == 0 {
		return Nil, fmt.Errorf("no prices found for date %s", date)
	}

	// Get thirty-day average
	avgPrice, err := GetThirtyDayAverage(ctx, date)
	if err != nil {
		return Nil, err
	}

	// Get day rating
	dayAvg := CalculateAverage(prices)
	dayRating := CalculateDayRating(dayAvg, avgPrice)

	return dayRating, nil
}

func GetDayAverage(ctx context.Context, date time.Time) (float64, error) {
	// Get the prices for the given day
	prices, err := GetDailyPrices(ctx, date)
	if err != nil {
		return 0, err
	}
	if len(prices) == 0 {
		return 0, fmt.Errorf("no prices found for date %s", date)
	}

	// Get day average
	dayAvg := CalculateAverage(prices)

	return dayAvg, nil
}

func GetCheapPeriods(ctx context.Context, date time.Time) ([][]Price, error) {
	// Get the prices for the given day
	prices, err := GetDailyPrices(ctx, date)
	if err != nil {
		return nil, err
	}
	if len(prices) == 0 {
		return nil, fmt.Errorf("no prices found for date %s", date)
	}

	// Get thirty-day average
	avgPrice, err := GetThirtyDayAverage(ctx, date)
	if err != nil {
		return nil, err
	}

	// Get cheap periods
	cheapPeriods := CalculateCheapPeriods(prices, avgPrice)

	return cheapPeriods, nil
}

func GetExpensivePeriods(ctx context.Context, date time.Time) ([][]Price, error) {
	// Get the prices for the given day
	prices, err := GetDailyPrices(ctx, date)
	if err != nil {
		return nil, err
	}
	if len(prices) == 0 {
		return nil, fmt.Errorf("no prices found for date %s", date)
	}

	// Get thirty-day average
	avgPrice, err := GetThirtyDayAverage(ctx, date)
	if err != nil {
		return nil, err
	}

	// Get expensive periods
	expensivePeriods := CalculateExpensivePeriods(prices, avgPrice)

	return expensivePeriods, nil
}
