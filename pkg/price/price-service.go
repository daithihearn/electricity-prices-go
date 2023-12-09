package price

import (
	"context"
	"electricity-prices/pkg/date"
	"electricity-prices/pkg/db"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
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
	Collection db.Collection[Price]
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

	client, err := db.GetMongoClient(ctx)
	if err != nil {
		log.Fatalf("Error getting mongo ree: %v", err)
	}

	// Insert the documents
	// Start a session for the transaction.
	session, err := client.StartSession()
	if err != nil {
		log.Fatalf("Error starting session: %v", err)
	}
	defer session.EndSession(ctx)

	// Define the work to be done in the transaction.
	txnErr := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		// Start the transaction
		err := session.StartTransaction()
		if err != nil {
			return err
		}

		err = r.Collection.InsertMany(ctx, prices)
		if err != nil {
			// If there'r an error, abort the transaction and return the error.
			session.AbortTransaction(sc)
			return err
		}

		// If everything went well, commit the transaction.
		err = session.CommitTransaction(sc)
		return err
	})

	if txnErr != nil {
		log.Fatalf("Transaction failed: %v", txnErr)
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
	start, end := date.ParseStartAndEndTimes(t, 30)

	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"dateTime": bson.M{
				"$gte": start,
				"$lte": end,
			},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id": nil,
			"averagePrice": bson.M{
				"$avg": "$price",
			},
		}}},
		{{Key: "$project", Value: bson.M{
			"_id":          0,
			"averagePrice": 1,
		}}},
	}

	cursor, err := r.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(cursor, ctx)

	var result bson.M
	if cursor.Next(ctx) {
		if err = cursor.Decode(&result); err != nil {
			return 0, err
		}
		if avgPrice, ok := result["averagePrice"].(float64); ok {
			return avgPrice, nil
		} else {
			return 0, fmt.Errorf("failed to convert average price to float64")
		}
	}

	return 0, fmt.Errorf("no results found")
}

// GetLatestPrice returns the latest price from the database
// It returns a boolean indicating if no price was found
// and an error if there was one
func (r *Receiver) GetLatestPrice(ctx context.Context) (Price, bool, error) {
	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{
		{{Key: "$sort", Value: bson.M{
			"dateTime": -1,
		}}},
		{{Key: "$limit", Value: 1}},
	}

	cursor, err := r.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return Price{}, false, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(cursor, ctx)

	var result Price
	if cursor.Next(ctx) {
		if err = cursor.Decode(&result); err != nil {
			return Price{}, false, err
		}
		return result, false, nil
	}

	return Price{}, true, nil
}
