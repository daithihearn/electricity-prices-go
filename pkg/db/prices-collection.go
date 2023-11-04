package db

import (
	"context"
	"electricity-prices/pkg/model"
	"electricity-prices/pkg/utils"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetPrice(now time.Time) (model.Price, error) {
	// Set to the start of the current hour
	hour := now.Truncate(time.Hour)

	// Get the prices for the given hour
	filter := bson.M{
		"dateTime": hour,
	}

	var price model.Price
	err := GetCollection().FindOne(context.Background(), filter).Decode(&price)

	return price, err
}

func SavePrices(prices []model.Price) error {
	collection := GetCollection()

	var documents []interface{}
	for _, price := range prices {
		documents = append(documents, price)
	}

	client, err := GetMongoClient()
	if err != nil {
		log.Fatalf("Error getting mongo client: %v", err)
	}

	// Insert the documents
	// Start a session for the transaction.
	session, err := client.StartSession()
	if err != nil {
		log.Fatalf("Error starting session: %v", err)
	}
	defer session.EndSession(context.Background())

	// Define the work to be done in the transaction.
	txnErr := mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
		// Start the transaction
		err := session.StartTransaction()
		if err != nil {
			return err
		}

		_, err = collection.InsertMany(context.Background(), documents)
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

func GetPrices(start time.Time, end time.Time) ([]model.Price, error) {

	collection := GetCollection()

	// Create a filter based on the date range
	filter := bson.M{
		"dateTime": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}

	findOptions := options.Find()
	var prices = make([]model.Price, 0)

	cur, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var price model.Price
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

func GetThirtyDayAverage(date time.Time) (float64, error) {
	start, end := utils.ParseStartAndEndTimes(date, 30)

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

	cursor, err := ExecutePipeline(pipeline)
	if err != nil {
		return 0, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(cursor, context.TODO())

	var result bson.M
	if cursor.Next(context.TODO()) {
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

func GetLatestPrice() (model.Price, error) {
	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{
		{{Key: "$sort", Value: bson.M{
			"dateTime": -1,
		}}},
		{{Key: "$limit", Value: 1}},
	}

	cursor, err := ExecutePipeline(pipeline)
	if err != nil {
		return model.Price{}, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(cursor, context.Background())

	var result model.Price
	if cursor.Next(context.Background()) {
		if err = cursor.Decode(&result); err != nil {
			return model.Price{}, err
		}
		return result, nil
	}

	return model.Price{}, fmt.Errorf("no results found")
}
