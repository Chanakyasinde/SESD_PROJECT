package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"inventory-management-system/backend/internal/domain/entities"
	"inventory-management-system/backend/internal/domain/repositories"
)

type StockLogRepository struct {
	collection *mongo.Collection
}

func NewStockLogRepository(collection *mongo.Collection) repositories.StockLogRepository {
	return &StockLogRepository{collection: collection}
}

func (r *StockLogRepository) Create(ctx context.Context, log *entities.StockLog) error {
	_, err := r.collection.InsertOne(ctx, log)
	return err
}
