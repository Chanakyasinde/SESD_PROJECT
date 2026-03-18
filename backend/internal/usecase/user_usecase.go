package usecase

import (
	"context"

	"inventory-management-system/backend/internal/domain/entities"
	"inventory-management-system/backend/internal/domain/repositories"
	"inventory-management-system/backend/internal/domain/usecases"
)

type UserUsecase struct {
	userRepo repositories.UserRepository
}

func NewUserUsecase(userRepo repositories.UserRepository) usecases.UserUsecase {
	return &UserUsecase{userRepo: userRepo}
}

func (u *UserUsecase) ListCustomers(ctx context.Context, page, pageSize int64) ([]entities.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 50
	}

	return u.userRepo.List(ctx, repositories.UserFilter{
		Role: entities.RoleCustomer,
		Pagination: entities.Pagination{
			Page:     page,
			PageSize: pageSize,
			SortBy:   "createdAt",
			SortDir:  -1,
		},
	})
}
