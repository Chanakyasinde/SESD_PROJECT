package mongo

import (
	"context"
	"errors"
	"fmt"

	"inventory-management-system/backend/internal/domain/entities"
	"inventory-management-system/backend/internal/domain/repositories"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(collection *mongo.Collection) repositories.ProductRepository {
	return &ProductRepository{collection: collection}
}

func (r *ProductRepository) Create(ctx context.Context, product *entities.Product) error {
	result, err := r.collection.InsertOne(ctx, product)
	if err != nil {
		return err
	}
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		product.ID = oid
	}
	return nil
}

func (r *ProductRepository) Update(ctx context.Context, product *entities.Product) error {
	filter := bson.M{"_id": product.ID, "version": product.Version}
	update := bson.M{
		"$set": bson.M{
			"name":              product.Name,
			"description":       product.Description,
			"category":          product.Category,
			"price":             product.Price,
			"stockQuantity":     product.StockQuantity,
			"lowStockThreshold": product.LowStockThreshold,
			"updatedAt":         product.UpdatedAt,
		},
		"$inc": bson.M{"version": 1},
	}
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return entities.ErrConflict
	}
	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return entities.ErrNotFound
	}
	return nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Product, error) {
	var product entities.Product
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, entities.ErrNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) List(ctx context.Context, filter repositories.ProductFilter) ([]entities.Product, int64, error) {
	query := bson.M{}
	if filter.Category != "" {
		query["category"] = filter.Category
	}
	if filter.Search != "" {
		query["$text"] = bson.M{"$search": filter.Search}
	}
	if filter.MinPrice != nil || filter.MaxPrice != nil {
		priceQuery := bson.M{}
		if filter.MinPrice != nil {
			priceQuery["$gte"] = *filter.MinPrice
		}
		if filter.MaxPrice != nil {
			priceQuery["$lte"] = *filter.MaxPrice
		}
		query["price"] = priceQuery
	}
	if filter.LowStock != nil && *filter.LowStock {
		query["$expr"] = bson.M{"$lte": []any{"$stockQuantity", "$lowStockThreshold"}}
	}

	count, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	sortDir := 1
	if filter.SortDir < 0 {
		sortDir = -1
	}
	sortBy := filter.SortBy
	if sortBy == "" {
		sortBy = "createdAt"
	}

	opts := options.Find().SetSkip((filter.Page - 1) * filter.PageSize).SetLimit(filter.PageSize).SetSort(bson.D{{Key: sortBy, Value: sortDir}})
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	products := make([]entities.Product, 0)
	for cursor.Next(ctx) {
		var product entities.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, 0, err
		}
		products = append(products, product)
	}
	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	return products, count, nil
}

func (r *ProductRepository) DecreaseStockWithVersion(ctx context.Context, id primitive.ObjectID, qty int64, version int64) error {
	if qty <= 0 {
		return fmt.Errorf("qty must be > 0")
	}

	filter := bson.M{"_id": id, "version": version, "stockQuantity": bson.M{"$gte": qty}}
	update := bson.M{"$inc": bson.M{"stockQuantity": -qty, "version": 1}}
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return entities.ErrConflict
	}
	return nil
}

func (r *ProductRepository) IncreaseStock(ctx context.Context, id primitive.ObjectID, qty int64) error {
	if qty <= 0 {
		return fmt.Errorf("qty must be > 0")
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"stockQuantity": qty, "version": 1}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return entities.ErrNotFound
	}
	return nil
}
