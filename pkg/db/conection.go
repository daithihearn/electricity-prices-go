package db

import (
	"context"
	"errors"
	"log"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientInstance      *mongo.Client
	clientInstanceError error
	mongoOnce           sync.Once
)

func GetMongoClient() (*mongo.Client, error) {
	mongoOnce.Do(func() {
		connectionString := os.Getenv("MONGODB_URI")
		if connectionString == "" {
			log.Fatal("MONGODB_URI must be set", errors.New("MONGODB_URI must be set"))
			return
		}

		// Define a command monitor
		//cmdMonitor := &event.CommandMonitor{
		//	Started: func(_ context.Context, evt *event.CommandStartedEvent) {
		//		log.Printf("Started command %s with data: %v", evt.CommandName, evt.Command)
		//	},
		//	Succeeded: func(_ context.Context, evt *event.CommandSucceededEvent) {
		//		log.Printf("Succeeded command %s with result: %v", evt.CommandName, evt.Reply)
		//	},
		//	Failed: func(_ context.Context, evt *event.CommandFailedEvent) {
		//		log.Printf("Failed command %s with error: %v", evt.CommandName, evt.Failure)
		//	},
		//}

		log.Println("Connected to MongoDB:", connectionString)

		// Set client options
		clientOptions := options.Client().ApplyURI(connectionString)

		// Connect to MongoDB
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			clientInstanceError = err
			return
		}

		// Check the connection
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			clientInstanceError = err
			return
		}

		clientInstance = client
	})
	return clientInstance, clientInstanceError
}

func GetCollection() *mongo.Collection {
	client, err := GetMongoClient()
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("electricity-prices").Collection("prices")

	return collection
}

func ExecutePipeline(pipeline mongo.Pipeline) (*mongo.Cursor, error) {
	collection := GetCollection()
	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

// CloseMongoConnection will close the MongoDB connection when the application exits.
func CloseMongoConnection() {
	if clientInstance != nil {
		err := clientInstance.Disconnect(context.TODO())
		if err != nil {
			log.Fatalf("Failed to close MongoDB connection: %v", err)
		}
	}
}
