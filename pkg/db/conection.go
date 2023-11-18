package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"sync"
)

var (
	clientInstance      *mongo.Client
	clientInstanceError error
	mongoOnce           sync.Once
)

// GetMongoClient initializes and returns a MongoDB client instance.
func GetMongoClient(ctx context.Context) (*mongo.Client, error) {
	mongoOnce.Do(func() {
		connectionString := os.Getenv("MONGODB_URI")
		if connectionString == "" {
			clientInstanceError = errors.New("MONGODB_URI must be set")
			log.Fatal(clientInstanceError)
			return
		}

		log.Println("Connecting to MongoDB:", connectionString)

		clientOptions := options.Client().ApplyURI(connectionString)

		var err error
		clientInstance, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			clientInstanceError = err
			return
		}

		err = clientInstance.Ping(ctx, nil)
		if err != nil {
			clientInstanceError = err
			return
		}
	})
	return clientInstance, clientInstanceError
}

func GetCollection(ctx context.Context) (*mongo.Collection, error) {
	client, err := GetMongoClient(ctx)
	if err != nil {
		return nil, err
	}

	db := client.Database("electricity-prices")
	collection := db.Collection("prices")

	return collection, nil
}

func ExecutePipeline(ctx context.Context, pipeline mongo.Pipeline) (*mongo.Cursor, error) {
	collection, err := GetCollection(ctx)
	if err != nil {
		return nil, err
	}

	return collection.Aggregate(ctx, pipeline)
}

func CloseMongoConnection(ctx context.Context) error {
	if clientInstance != nil {
		return clientInstance.Disconnect(ctx)
	}
	return nil
}

// In your main function or where you call these functions, you should create a context:
// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// defer cancel()
