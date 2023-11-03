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
