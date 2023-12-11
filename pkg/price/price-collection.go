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

type Collection interface {
	db.Collection[Price]
	GetThirtyDayAverage(ctx context.Context, t time.Time) (float64, error)
	GetLatestPrice(ctx context.Context) (Price, bool, error)
}

type ColReceiver struct {
	Col *mongo.Collection
}

func (r ColReceiver) FindOne(ctx context.Context, filter interface{}) (Price, error) {
	var p Price
	err := r.Col.FindOne(ctx, filter).Decode(&p)

	if err != nil {
		return Price{}, err
	}

	return p, err
}

func (r ColReceiver) Find(ctx context.Context, filter interface{}) ([]Price, error) {
	cur, err := r.Col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(cur, ctx)

	var prices = make([]Price, 0)

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

func (r ColReceiver) InsertMany(ctx context.Context, documents []Price) error {
	client, err := db.GetMongoClient(ctx)
	if err != nil {
		log.Fatalf("Error getting mongo client: %v", err)
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
		var documentsInterface []interface{}
		for _, doc := range documents {
			documentsInterface = append(documentsInterface, doc)
		}

		_, err = r.Col.InsertMany(ctx, documentsInterface)
		if err != nil {
			return err
		}

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

func (r ColReceiver) Aggregate(ctx context.Context, pipeline interface{}) (*mongo.Cursor, error) {
	cursor, err := r.Col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	return cursor, nil
}

func (r ColReceiver) GetThirtyDayAverage(ctx context.Context, t time.Time) (float64, error) {
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

	cursor, err := r.Aggregate(ctx, pipeline)
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

func (r ColReceiver) GetLatestPrice(ctx context.Context) (Price, bool, error) {
	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{
		{{Key: "$sort", Value: bson.M{
			"dateTime": -1,
		}}},
		{{Key: "$limit", Value: 1}},
	}

	cursor, err := r.Aggregate(ctx, pipeline)
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
