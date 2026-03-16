package repositories

import (
	"context"

	"inventory-management-system/backend/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserFilter struct {
	Role string
	entities.Pagination
}

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*entities.User, error)
	List(ctx context.Context, filter UserFilter) ([]entities.User, int64, error)
}
