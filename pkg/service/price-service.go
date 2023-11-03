package service

import (
	"electricity-prices/pkg/db"
	"electricity-prices/pkg/model"
	"electricity-prices/pkg/utils"
	"fmt"
	"time"
)

func GetDailyPrices(date time.Time) ([]model.Price, error) {
	start, end := utils.ParseStartAndEndTimes(date, 0)

	prices, err := db.GetPrices(start, end)

	if err != nil {
		return nil, err
	}

	return prices, nil
}

func GetDailyAverages(date time.Time, numberOfDays int) ([]model.DailyAverage, error) {

	xDaysAgo := date.AddDate(0, 0, -numberOfDays)
	nextDay := time.Date(date.Year(), date.Month(), date.Day()+1, 0, 0, 0, 0, date.Location())

	// Subtract one second to get the last second of the current day
	today := nextDay.Add(-time.Second)

	prices, err := db.GetPrices(xDaysAgo, today)

	if err != nil {
		return nil, err
	}

	averages := utils.CalculateDailyAverages(prices)

	return averages, nil

}

func GetDailyInfo(date time.Time) (model.DailyPriceInfo, error) {
	// Get the prices for the given day
	prices, err := GetDailyPrices(date)
	if err != nil {
		return model.DailyPriceInfo{}, err
	}

	// Get thirty-day average
	avgPrice, err := db.GetThirtyDayAverage(date)
	if err != nil {
		return model.DailyPriceInfo{}, err
	}

	// Get day rating
	dayAvg := utils.CalculateAverage(prices)
	dayRating := utils.CalculateDayRating(dayAvg, avgPrice)

	// Get cheap periods
	cheapPeriods := utils.CalculateCheapPeriods(prices, avgPrice)

	// Get expensive periods
	expensivePeriods := utils.CalculateExpensivePeriods(prices, avgPrice)

	return model.DailyPriceInfo{
		Prices:           prices,
		ThirtyDayAverage: avgPrice,
		DayRating:        dayRating,
		DayAverage:       dayAvg,
		CheapPeriods:     cheapPeriods,
		ExpensivePeriods: expensivePeriods,
	}, nil
}

func GetDayRating(date time.Time) (model.DayRating, error) {
	// Get the prices for the given day
	prices, err := GetDailyPrices(date)
	if err != nil {
		return model.Nil, err
	}
	if len(prices) == 0 {
		return model.Nil, fmt.Errorf("no prices found for date %s", date)
	}

	// Get thirty-day average
	avgPrice, err := db.GetThirtyDayAverage(date)
	if err != nil {
		return model.Nil, err
	}

	// Get day rating
	dayAvg := utils.CalculateAverage(prices)
	dayRating := utils.CalculateDayRating(dayAvg, avgPrice)

	return dayRating, nil
}

func GetDayAverage(date time.Time) (float64, error) {
	// Get the prices for the given day
	prices, err := GetDailyPrices(date)
	if err != nil {
		return 0, err
	}
	if len(prices) == 0 {
		return 0, fmt.Errorf("no prices found for date %s", date)
	}

	// Get day average
	dayAvg := utils.CalculateAverage(prices)

	return dayAvg, nil
}

func GetCheapPeriods(date time.Time) ([][]model.Price, error) {
	// Get the prices for the given day
	prices, err := GetDailyPrices(date)
	if err != nil {
		return nil, err
	}
	if len(prices) == 0 {
		return nil, fmt.Errorf("no prices found for date %s", date)
	}

	// Get thirty-day average
	avgPrice, err := db.GetThirtyDayAverage(date)
	if err != nil {
		return nil, err
	}

	// Get cheap periods
	cheapPeriods := utils.CalculateCheapPeriods(prices, avgPrice)

	return cheapPeriods, nil
}

func GetExpensivePeriods(date time.Time) ([][]model.Price, error) {
	// Get the prices for the given day
	prices, err := GetDailyPrices(date)
	if err != nil {
		return nil, err
	}
	if len(prices) == 0 {
		return nil, fmt.Errorf("no prices found for date %s", date)
	}

	// Get thirty-day average
	avgPrice, err := db.GetThirtyDayAverage(date)
	if err != nil {
		return nil, err
	}

	// Get expensive periods
	expensivePeriods := utils.CalculateExpensivePeriods(prices, avgPrice)

	return expensivePeriods, nil
}

func GetPrice(date time.Time) (model.Price, error) {
	return db.GetPrice(date)
}

func GetThirtyDayAverage(date time.Time) (float64, error) {
	return db.GetThirtyDayAverage(date)
}
