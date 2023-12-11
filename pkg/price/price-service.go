package price

import (
	"context"
	"electricity-prices/pkg/date"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type Service interface {
	GetPrice(ctx context.Context, t time.Time) (Price, error)
	GetPrices(ctx context.Context, start time.Time, end time.Time) ([]Price, error)
	SavePrices(ctx context.Context, prices []Price) error
	GetDailyPrices(ctx context.Context, t time.Time) ([]Price, error)
	GetDailyAverages(ctx context.Context, t time.Time, numberOfDays int) ([]DailyAverage, error)
	GetDailyInfo(ctx context.Context, t time.Time) (DailyPriceInfo, error)
	GetDayRating(ctx context.Context, t time.Time) (DayRating, error)
	GetDayAverage(ctx context.Context, t time.Time) (float64, error)
	GetCheapPeriods(ctx context.Context, t time.Time) ([][]Price, error)
	GetExpensivePeriods(ctx context.Context, t time.Time) ([][]Price, error)
	GetThirtyDayAverage(ctx context.Context, t time.Time) (float64, error)
	GetLatestPrice(ctx context.Context) (Price, bool, error)
}

type Receiver struct {
	Collection Collection
}

func (r *Receiver) GetPrice(ctx context.Context, t time.Time) (Price, error) {
	// Set to the start of the current hour
	hour := t.Truncate(time.Hour)

	// Get the prices for the given hour
	filter := bson.M{
		"dateTime": hour,
	}

	return r.Collection.FindOne(ctx, filter)
}

func (r *Receiver) GetPrices(ctx context.Context, start time.Time, end time.Time) ([]Price, error) {

	// Create a filter based on the date range
	filter := bson.M{
		"dateTime": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}

	return r.Collection.Find(ctx, filter)
}

func (r *Receiver) SavePrices(ctx context.Context, prices []Price) error {
	err := r.Collection.InsertMany(ctx, prices)
	if err != nil {
		return err
	}
	return nil
}

func (r *Receiver) GetDailyPrices(ctx context.Context, t time.Time) ([]Price, error) {
	start, end := date.ParseStartAndEndTimes(t, 1)

	prices, err := r.GetPrices(ctx, start, end)

	if err != nil {
		return nil, err
	}

	return prices, nil
}

func (r *Receiver) GetDailyAverages(ctx context.Context, t time.Time, numberOfDays int) ([]DailyAverage, error) {

	xDaysAgo := t.AddDate(0, 0, -numberOfDays)
	nextDay := time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, t.Location())

	// Subtract one second to get the last second of the current day
	today := nextDay.Add(-time.Second)

	prices, err := r.GetPrices(ctx, xDaysAgo, today)

	if err != nil {
		return nil, err
	}

	averages := CalculateDailyAverages(prices)

	return averages, nil

}

func (r *Receiver) GetDailyInfo(ctx context.Context, t time.Time) (DailyPriceInfo, error) {
	// Get the prices for the given day
	prices, err := r.GetDailyPrices(ctx, t)
	if err != nil {
		return DailyPriceInfo{}, err
	}

	// Get thirty-day average
	avgPrice, err := r.GetThirtyDayAverage(ctx, t)
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

func (r *Receiver) GetDayRating(ctx context.Context, t time.Time) (DayRating, error) {
	// Get the prices for the given day
	prices, err := r.GetDailyPrices(ctx, t)
	if err != nil {
		return Nil, err
	}
	if len(prices) == 0 {
		return Nil, fmt.Errorf("no prices found for t %s", t)
	}

	// Get thirty-day average
	avgPrice, err := r.GetThirtyDayAverage(ctx, t)
	if err != nil {
		return Nil, err
	}

	// Get day rating
	dayAvg := CalculateAverage(prices)
	dayRating := CalculateDayRating(dayAvg, avgPrice)

	return dayRating, nil
}

func (r *Receiver) GetDayAverage(ctx context.Context, t time.Time) (float64, error) {
	// Get the prices for the given day
	prices, err := r.GetDailyPrices(ctx, t)
	if err != nil {
		return 0, err
	}
	if len(prices) == 0 {
		return 0, fmt.Errorf("no prices found for t %s", t)
	}

	// Get day average
	dayAvg := CalculateAverage(prices)

	return dayAvg, nil
}

func (r *Receiver) GetCheapPeriods(ctx context.Context, t time.Time) ([][]Price, error) {
	// Get the prices for the given day
	prices, err := r.GetDailyPrices(ctx, t)
	if err != nil {
		return nil, err
	}
	if len(prices) == 0 {
		return nil, fmt.Errorf("no prices found for t %s", t)
	}

	// Get thirty-day average
	avgPrice, err := r.GetThirtyDayAverage(ctx, t)
	if err != nil {
		return nil, err
	}

	// Get cheap periods
	cheapPeriods := CalculateCheapPeriods(prices, avgPrice)

	return cheapPeriods, nil
}

func (r *Receiver) GetExpensivePeriods(ctx context.Context, t time.Time) ([][]Price, error) {
	// Get the prices for the given day
	prices, err := r.GetDailyPrices(ctx, t)
	if err != nil {
		return nil, err
	}
	if len(prices) == 0 {
		return nil, fmt.Errorf("no prices found for t %s", t)
	}

	// Get thirty-day average
	avgPrice, err := r.GetThirtyDayAverage(ctx, t)
	if err != nil {
		return nil, err
	}

	// Get expensive periods
	expensivePeriods := CalculateExpensivePeriods(prices, avgPrice)

	return expensivePeriods, nil
}

func (r *Receiver) GetThirtyDayAverage(ctx context.Context, t time.Time) (float64, error) {
	return r.Collection.GetThirtyDayAverage(ctx, t)
}

// GetLatestPrice returns the latest price from the database
// It returns a boolean indicating if no price was found
// and an error if there was one
func (r *Receiver) GetLatestPrice(ctx context.Context) (Price, bool, error) {
	return r.Collection.GetLatestPrice(ctx)
}
