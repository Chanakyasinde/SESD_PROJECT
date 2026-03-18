package usecases

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"inventory-management-system/backend/internal/domain/entities"
	"inventory-management-system/backend/internal/domain/repositories"
)

type CreateOrderItemInput struct {
	ProductID string `json:"productId"`
	Quantity  int64  `json:"quantity"`
}

type CreateOrderInput struct {
	CustomerID string                 `json:"customerId"`
	Items      []CreateOrderItemInput `json:"items"`
}

type OrderUsecase interface {
	Create(ctx context.Context, actorID primitive.ObjectID, input CreateOrderInput) (*entities.Order, error)
	UpdateStatus(ctx context.Context, actorID primitive.ObjectID, id primitive.ObjectID, status string) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Order, error)
	List(ctx context.Context, filter repositories.OrderFilter) ([]entities.Order, int64, error)
}
