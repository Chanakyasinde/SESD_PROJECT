package usecases

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"inventory-management-system/backend/internal/domain/entities"
	"inventory-management-system/backend/internal/domain/repositories"
)

type ProductUsecase interface {
	Create(ctx context.Context, input entities.Product) (*entities.Product, error)
	Update(ctx context.Context, id primitive.ObjectID, input entities.Product) (*entities.Product, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Product, error)
	List(ctx context.Context, filter repositories.ProductFilter) ([]entities.Product, int64, error)
}
