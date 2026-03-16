package mongo

import (
	"context"
	"errors"
	"time"

	"inventory-management-system/backend/internal/domain/entities"
	"inventory-management-system/backend/internal/domain/repositories"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderRepository struct {
	collection *mongo.Collection
}

func NewOrderRepository(collection *mongo.Collection) repositories.OrderRepository {
	return &OrderRepository{collection: collection}
}

func (r *OrderRepository) Create(ctx context.Context, order *entities.Order) error {
	result, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		return err
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("failed to parse inserted order id")
	}
	order.ID = oid
	return nil
}

func (r *OrderRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Order, error) {
	var order entities.Order
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, entities.ErrNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status, "updatedAt": time.Now()}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return entities.ErrNotFound
	}
	return nil
}

func (r *OrderRepository) List(ctx context.Context, filter repositories.OrderFilter) ([]entities.Order, int64, error) {
	query := bson.M{}
	if filter.Status != "" {
		query["status"] = filter.Status
	}
	if filter.CustomerID != nil {
		query["customerId"] = *filter.CustomerID
	}

	count, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	sortBy := filter.SortBy
	if sortBy == "" {
		sortBy = "createdAt"
	}
	sortDir := 1
	if filter.SortDir < 0 {
		sortDir = -1
	}
	opts := options.Find().SetSkip((filter.Page - 1) * filter.PageSize).SetLimit(filter.PageSize).SetSort(bson.D{{Key: sortBy, Value: sortDir}})

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	orders := make([]entities.Order, 0)
	for cursor.Next(ctx) {
		var order entities.Order
		if err := cursor.Decode(&order); err != nil {
			return nil, 0, err
		}
		orders = append(orders, order)
	}
	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	return orders, count, nil
}
