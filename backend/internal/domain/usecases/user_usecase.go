package usecases

import (
	"context"

	"inventory-management-system/backend/internal/domain/entities"
)

type UserUsecase interface {
	ListCustomers(ctx context.Context, page, pageSize int64) ([]entities.User, int64, error)
}
