package price

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type PriceCollection struct {
	Col *mongo.Collection
}

func (c PriceCollection) FindOne(ctx context.Context, filter interface{}) (Price, error) {
	var p Price
	err := c.Col.FindOne(ctx, filter).Decode(&p)

	if err != nil {
		return Price{}, err
	}

	return p, err
}

func (c PriceCollection) Find(ctx context.Context, filter interface{}) ([]Price, error) {
	cur, err := c.Col.Find(ctx, filter)
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

func (c PriceCollection) InsertMany(ctx context.Context, documents []Price) error {
	var documentsInterface []interface{}
	for _, doc := range documents {
		documentsInterface = append(documentsInterface, doc)
	}

	_, err := c.Col.InsertMany(ctx, documentsInterface)
	if err != nil {
		return err
	}

	return nil
}

func (c PriceCollection) Aggregate(ctx context.Context, pipeline interface{}) (*mongo.Cursor, error) {
	cursor, err := c.Col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	return cursor, nil
}
