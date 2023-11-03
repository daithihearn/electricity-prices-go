package utils

import (
	"electricity-prices/pkg/model"
	"fmt"
	"strings"
	"time"
)

const ratingVariance = 0.02
const varianceDivisor = 2.0

func CalculateDailyAverages(prices []model.Price) []model.DailyAverage {
	averages := make([]model.DailyAverage, 0, len(prices)/24)

	// Group the prices by date
	dateMap := make(map[string][]model.Price)
	for _, price := range prices {
		dateOnly := ParseToLocalDay(price.DateTime)
		dateMap[dateOnly] = append(dateMap[dateOnly], price)
	}

	// Calculate the average for each day
	for date, prices := range dateMap {
		average := CalculateAverage(prices)
		averages = append(averages, model.DailyAverage{Date: date, Average: average})
	}

	// Sort the averages by date
	sortDailyAverages(averages)

	return averages
}

func CalculateAverage(prices []model.Price) float64 {
	if len(prices) == 0 {
		return 0.0
	}
	var total float64
	for _, price := range prices {
		total += price.Price
	}
	return total / float64(len(prices))
}

// CalculateDayRating
// Calculate the rating for a day based on the daily average and the thirty-day average.
func CalculateDayRating(dayAvg float64, thirtyDayAvg float64) model.DayRating {

	variance := dayAvg - thirtyDayAvg

	if variance < -ratingVariance {
		return model.Good
	} else if variance > ratingVariance {
		return model.Bad
	}
	return model.Normal
}

// CalculateCheapPeriods
// Get the cheap periods from a slice of prices.
func CalculateCheapPeriods(prices []model.Price, thirtyDayAvg float64) [][]model.Price {
	if len(prices) == 0 {
		return [][]model.Price{}
	}
	variance := calculateCheapVariance(prices, thirtyDayAvg)
	minP := getMinPrice(prices)

	// Filter prices that are below the variance
	var cheapPeriods []model.Price
	for _, price := range prices {
		if price.Price <= minP+variance {
			cheapPeriods = append(cheapPeriods, price)
		}
	}

	// Group consecutive prices into periods
	return groupPrices(cheapPeriods)
}

// CalculateExpensivePeriods
// Get the expensive periods from a slice of prices.
func CalculateExpensivePeriods(prices []model.Price, thirtyDayAvg float64) [][]model.Price {
	if len(prices) == 0 {
		return [][]model.Price{}
	}
	variance := calculateExpensiveVariance(prices, thirtyDayAvg)
	maxP := getMaxPrice(prices)

	// If the most expensive price would be considered cheap, return empty list
	if maxP <= thirtyDayAvg-ratingVariance {
		return [][]model.Price{}
	}

	// Filter prices that are above the variance
	var expensivePeriods []model.Price
	for _, price := range prices {
		if maxP-price.Price <= variance {
			expensivePeriods = append(expensivePeriods, price)
		}
	}

	// Group consecutive prices into periods
	return groupPrices(expensivePeriods)
}

func sortDailyAverages(averages []model.DailyAverage) {
	// Sort the averages by date
	for i := 0; i < len(averages); i++ {
		for j := i + 1; j < len(averages); j++ {
			if averages[i].Date > averages[j].Date {
				temp := averages[i]
				averages[i] = averages[j]
				averages[j] = temp
			}
		}
	}
}

// calculateCombinedAverage
// Calculate the combination between the thirty-day average and the daily average.
// A weighted average is used with the daily average being weighted twice as much as the thirty-day average.
func calculateCombinedAverage(dayAvg float64, thirtyDayAvg float64) float64 {
	if dayAvg == 0.0 && thirtyDayAvg == 0.0 {
		return 0.0
	}
	return (dayAvg*2 + thirtyDayAvg) / 3
}

// getMinPrice
// Get the minimum price from a slice of prices.
func getMinPrice(prices []model.Price) float64 {
	if len(prices) == 0 {
		return 0.0
	}
	minP := prices[0].Price
	for _, price := range prices {
		if price.Price < minP {
			minP = price.Price
		}
	}
	return minP
}

// getMaxPrice
// Get the maximum price from a slice of prices.
func getMaxPrice(prices []model.Price) float64 {
	if len(prices) == 0 {
		return 0.0
	}
	maxP := prices[0].Price
	for _, price := range prices {
		if price.Price > maxP {
			maxP = price.Price
		}
	}
	return maxP
}

// getMinAndMaxPrices
// Get the maximum and minimum prices from a slice of prices.
func getMinAndMaxPrices(prices []model.Price) (float64, float64) {
	if len(prices) == 0 {
		return 0.0, 0.0
	}
	minP := prices[0].Price
	maxP := prices[0].Price
	for _, price := range prices {
		if price.Price < minP {
			minP = price.Price
		}
		if price.Price > maxP {
			maxP = price.Price
		}
	}
	return minP, maxP
}

// calculateMinVariance
// Calculate the minimum allowable variance. Calculated as a 6th of the distance between the most expensive and cheapest prices.
func calculateMinVariance(minP, maxP float64) float64 {
	return (maxP - minP) / 6
}

// calculateMaxVariance
// Calculate the maximum allowable variance. Calculated as a 3rd of the distance between the most expensive and cheapest prices.
func calculateMaxVariance(minP, maxP float64) float64 {
	return (maxP - minP) / 3
}

// calculateCheapVariance
// Calculate the variance for a cheap periods.
// MAX(MIN((dayAvg - cheapestPrice) / varianceDivisor, maxVariance), minVariance)
func calculateCheapVariance(prices []model.Price, thirtyDayAvg float64) float64 {
	minP, maxP := getMinAndMaxPrices(prices)
	minVariance := calculateMinVariance(minP, maxP)
	maxVariance := calculateMaxVariance(minP, maxP)
	dayAvg := CalculateAverage(prices)
	combAvg := calculateCombinedAverage(dayAvg, thirtyDayAvg)
	variance := (combAvg - minP) / varianceDivisor
	if variance < minVariance {
		variance = minVariance
	} else if variance > maxVariance {
		variance = maxVariance
	}
	return variance
}

// calculateExpensiveVariance
// Calculate the variance for an expensive periods.
// MAX(MIN((mostExpensivePrice - dayAvg) / varianceDivisor, maxVariance), minVariance)
func calculateExpensiveVariance(prices []model.Price, thirtyDayAvg float64) float64 {
	minP, maxP := getMinAndMaxPrices(prices)
	minVariance := calculateMinVariance(minP, maxP)
	maxVariance := calculateMaxVariance(minP, maxP)
	dayAvg := CalculateAverage(prices)
	combAvg := calculateCombinedAverage(dayAvg, thirtyDayAvg)
	variance := (maxP - combAvg) / varianceDivisor
	if variance < minVariance {
		variance = minVariance
	} else if variance > maxVariance {
		variance = maxVariance
	}
	return variance
}

// groupPrices
// Group consecutive prices into periods.
func groupPrices(cheapPrices []model.Price) [][]model.Price {
	var result [][]model.Price
	i := 0
	for i < len(cheapPrices) {
		cheapPeriod := []model.Price{cheapPrices[i]}
		j := i + 1
		for j < len(cheapPrices) && cheapPrices[j].DateTime.Sub(cheapPeriod[len(cheapPeriod)-1].DateTime) == time.Hour {
			cheapPeriod = append(cheapPeriod, cheapPrices[j])
			j++
		}
		result = append(result, cheapPeriod)
		i = j
	}
	return result
}

// GetNextPeriod
// Given the provided date and price periods, return the next period.
// Also return whether the next period has started yet or not
func GetNextPeriod(prices [][]model.Price, date time.Time) ([]model.Price, bool) {

	// Get the next cheap period
	for _, period := range prices {
		if len(period) == 0 {
			continue
		}
		if period[0].DateTime.After(date) {
			return period, false
		} else if period[len(period)-1].DateTime.Add(59 * time.Minute).After(date) {
			return period, true
		}
	}

	return nil, false
}

// FormatPrice
// Format a price to a string with 2 decimal places.
func FormatPrice(price float64) string {
	p := fmt.Sprintf("%.1f", price*100)
	if strings.HasSuffix(p, ".0") {
		return p[:len(p)-2]
	}
	return p
}
