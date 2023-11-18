package price

import (
	"context"
	"electricity-prices/pkg/date"
	"electricity-prices/pkg/db"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetPrice(ctx context.Context, now time.Time) (Price, error) {
	// Set to the start of the current hour
	hour := now.Truncate(time.Hour)

	// Get the prices for the given hour
	filter := bson.M{
		"dateTime": hour,
	}

	var price Price
	col, err := db.GetCollection(ctx)
	if err != nil {
		return Price{}, err
	}
	err = col.FindOne(ctx, filter).Decode(&price)
	if err != nil {
		return Price{}, err
	}

	return price, err
}

func SavePrices(ctx context.Context, prices []Price) error {
	col, err := db.GetCollection(ctx)

	if err != nil {
		return err
	}

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

		_, err = col.InsertMany(ctx, documents)
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

func getPrices(ctx context.Context, start time.Time, end time.Time) ([]Price, error) {

	col, err := db.GetCollection(ctx)

	if err != nil {
		return nil, err
	}

	// Create a filter based on the date range
	filter := bson.M{
		"dateTime": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}

	findOptions := options.Find()
	var prices = make([]Price, 0)

	cur, err := col.Find(ctx, filter, findOptions)
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

func GetThirtyDayAverage(ctx context.Context, t time.Time) (float64, error) {
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

	cursor, err := db.ExecutePipeline(ctx, pipeline)
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

func GetLatestPrice(ctx context.Context) (Price, error) {
	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{
		{{Key: "$sort", Value: bson.M{
			"dateTime": -1,
		}}},
		{{Key: "$limit", Value: 1}},
	}

	cursor, err := db.ExecutePipeline(ctx, pipeline)
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
