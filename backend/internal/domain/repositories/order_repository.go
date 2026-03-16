package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"inventory-management-system/backend/internal/domain/entities"
)

type OrderFilter struct {
	Status     string
	CustomerID *primitive.ObjectID
	entities.Pagination
}

type OrderRepository interface {
	Create(ctx context.Context, order *entities.Order) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Order, error)
	UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error
	List(ctx context.Context, filter OrderFilter) ([]entities.Order, int64, error)
}
