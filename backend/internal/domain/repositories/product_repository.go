package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"inventory-management-system/backend/internal/domain/entities"
)

type ProductFilter struct {
	Category string
	Search   string
	MinPrice *float64
	MaxPrice *float64
	LowStock *bool
	entities.Pagination
}

type ProductRepository interface {
	Create(ctx context.Context, product *entities.Product) error
	Update(ctx context.Context, product *entities.Product) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Product, error)
	List(ctx context.Context, filter ProductFilter) ([]entities.Product, int64, error)
	DecreaseStockWithVersion(ctx context.Context, id primitive.ObjectID, qty int64, version int64) error
	IncreaseStock(ctx context.Context, id primitive.ObjectID, qty int64) error
}
