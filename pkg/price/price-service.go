package price

import (
	"context"
	"electricity-prices/pkg/date"
	"electricity-prices/pkg/db"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Service struct {
	Collection db.Collection
}

func (s *Service) GetPrice(ctx context.Context, now time.Time) (Price, error) {
	// Set to the start of the current hour
	hour := now.Truncate(time.Hour)

	// Get the prices for the given hour
	filter := bson.M{
		"dateTime": hour,
	}

	var price Price
	err := s.Collection.FindOne(ctx, filter).Decode(&price)

	if err != nil {
		return Price{}, err
	}

	return price, err
}

func (s *Service) GetPrices(ctx context.Context, start time.Time, end time.Time) ([]Price, error) {

	// Create a filter based on the date range
	filter := bson.M{
		"dateTime": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}

	findOptions := options.Find()
	var prices = make([]Price, 0)

	cur, err := s.Collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(cur, ctx)

	for cur.Next(ctx) {
		var price Price
		err := cur.Decode(&price)
		if err != nil {
			log.Println("Error decoding price:", err)
			continue
		}
		prices = append(prices, price)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return prices, nil
}

func (s *Service) SavePrices(ctx context.Context, prices []Price) error {

	var documents []interface{}
	for _, price := range prices {
		documents = append(documents, price)
	}

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

		_, err = s.Collection.InsertMany(ctx, documents)
		if err != nil {
			// If there's an error, abort the transaction and return the error.
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

func (s *Service) GetDailyPrices(ctx context.Context, t time.Time) ([]Price, error) {
	start, end := date.ParseStartAndEndTimes(t, 1)

	prices, err := s.GetPrices(ctx, start, end)

	if err != nil {
		return nil, err
	}

	return prices, nil
}

func (s *Service) GetDailyAverages(ctx context.Context, date time.Time, numberOfDays int) ([]DailyAverage, error) {

	xDaysAgo := date.AddDate(0, 0, -numberOfDays)
	nextDay := time.Date(date.Year(), date.Month(), date.Day()+1, 0, 0, 0, 0, date.Location())

	// Subtract one second to get the last second of the current day
	today := nextDay.Add(-time.Second)

	prices, err := s.GetPrices(ctx, xDaysAgo, today)

	if err != nil {
		return nil, err
	}

	averages := CalculateDailyAverages(prices)

	return averages, nil

}

func (s *Service) GetDailyInfo(ctx context.Context, date time.Time) (DailyPriceInfo, error) {
	// Get the prices for the given day
	prices, err := s.GetDailyPrices(ctx, date)
	if err != nil {
		return DailyPriceInfo{}, err
	}

	// Get thirty-day average
	avgPrice, err := s.GetThirtyDayAverage(ctx, date)
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

func (s *Service) GetDayRating(ctx context.Context, date time.Time) (DayRating, error) {
	// Get the prices for the given day
	prices, err := s.GetDailyPrices(ctx, date)
	if err != nil {
		return Nil, err
	}
	if len(prices) == 0 {
		return Nil, fmt.Errorf("no prices found for date %s", date)
	}

	// Get thirty-day average
	avgPrice, err := s.GetThirtyDayAverage(ctx, date)
	if err != nil {
		return Nil, err
	}

	// Get day rating
	dayAvg := CalculateAverage(prices)
	dayRating := CalculateDayRating(dayAvg, avgPrice)

	return dayRating, nil
}

func (s *Service) GetDayAverage(ctx context.Context, date time.Time) (float64, error) {
	// Get the prices for the given day
	prices, err := s.GetDailyPrices(ctx, date)
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

func (s *Service) GetCheapPeriods(ctx context.Context, date time.Time) ([][]Price, error) {
	// Get the prices for the given day
	prices, err := s.GetDailyPrices(ctx, date)
	if err != nil {
		return nil, err
	}
	if len(prices) == 0 {
		return nil, fmt.Errorf("no prices found for date %s", date)
	}

	// Get thirty-day average
	avgPrice, err := s.GetThirtyDayAverage(ctx, date)
	if err != nil {
		return nil, err
	}

	// Get cheap periods
	cheapPeriods := CalculateCheapPeriods(prices, avgPrice)

	return cheapPeriods, nil
}

func (s *Service) GetExpensivePeriods(ctx context.Context, date time.Time) ([][]Price, error) {
	// Get the prices for the given day
	prices, err := s.GetDailyPrices(ctx, date)
	if err != nil {
		return nil, err
	}
	if len(prices) == 0 {
		return nil, fmt.Errorf("no prices found for date %s", date)
	}

	// Get thirty-day average
	avgPrice, err := s.GetThirtyDayAverage(ctx, date)
	if err != nil {
		return nil, err
	}

	// Get expensive periods
	expensivePeriods := CalculateExpensivePeriods(prices, avgPrice)

	return expensivePeriods, nil
}

func (s *Service) GetThirtyDayAverage(ctx context.Context, t time.Time) (float64, error) {
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

	cursor, err := s.Collection.Aggregate(ctx, pipeline)
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

func (s *Service) GetLatestPrice(ctx context.Context) (Price, error) {
	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{
		{{Key: "$sort", Value: bson.M{
			"dateTime": -1,
		}}},
		{{Key: "$limit", Value: 1}},
	}

	cursor, err := s.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return Price{}, err
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
			return Price{}, err
		}
		return result, nil
	}

	return Price{}, fmt.Errorf("no results found")
}
