package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Collection[T any] interface {
	FindOne(ctx context.Context, filter interface{}) (T, error)
	Find(ctx context.Context, filter interface{}) ([]T, error)
	InsertMany(ctx context.Context, documents []T) error
	Aggregate(ctx context.Context, pipeline interface{}) (*mongo.Cursor, error)
}
