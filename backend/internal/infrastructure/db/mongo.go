package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCollections struct {
	Users     *mongo.Collection
	Products  *mongo.Collection
	Orders    *mongo.Collection
	StockLogs *mongo.Collection
}

func Connect(ctx context.Context, uri, dbName string) (*mongo.Client, *MongoCollections, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, fmt.Errorf("connect mongo: %w", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, nil, fmt.Errorf("ping mongo: %w", err)
	}

	db := client.Database(dbName)
	collections := &MongoCollections{
		Users:     db.Collection("users"),
		Products:  db.Collection("products"),
		Orders:    db.Collection("orders"),
		StockLogs: db.Collection("stock_logs"),
	}

	if err = ensureIndexes(ctx, collections); err != nil {
		return nil, nil, err
	}

	return client, collections, nil
}

func ensureIndexes(ctx context.Context, c *MongoCollections) error {
	userIndexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)},
	}
	if _, err := c.Users.Indexes().CreateMany(ctx, userIndexes); err != nil {
		return fmt.Errorf("create users indexes: %w", err)
	}

	productIndexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "name", Value: "text"}, {Key: "description", Value: "text"}}},
		{Keys: bson.D{{Key: "category", Value: 1}}},
		{Keys: bson.D{{Key: "stockQuantity", Value: 1}}},
	}
	if _, err := c.Products.Indexes().CreateMany(ctx, productIndexes); err != nil {
		return fmt.Errorf("create products indexes: %w", err)
	}

	orderIndexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "customerId", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "createdAt", Value: -1}}},
	}
	if _, err := c.Orders.Indexes().CreateMany(ctx, orderIndexes); err != nil {
		return fmt.Errorf("create orders indexes: %w", err)
	}

	stockLogIndexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "productId", Value: 1}}},
		{Keys: bson.D{{Key: "orderId", Value: 1}}},
		{Keys: bson.D{{Key: "createdAt", Value: -1}}},
	}
	if _, err := c.StockLogs.Indexes().CreateMany(ctx, stockLogIndexes); err != nil {
		return fmt.Errorf("create stock_logs indexes: %w", err)
	}

	return nil
}
