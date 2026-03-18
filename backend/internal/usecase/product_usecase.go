package usecase

import (
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"inventory-management-system/backend/internal/domain/entities"
	"inventory-management-system/backend/internal/domain/repositories"
	"inventory-management-system/backend/internal/domain/usecases"
)

type ProductUsecase struct {
	productRepo repositories.ProductRepository
}

func NewProductUsecase(productRepo repositories.ProductRepository) usecases.ProductUsecase {
	return &ProductUsecase{productRepo: productRepo}
}

func (u *ProductUsecase) Create(ctx context.Context, input entities.Product) (*entities.Product, error) {
	if strings.TrimSpace(input.Name) == "" || input.Price <= 0 || input.StockQuantity < 0 {
		return nil, entities.ErrInvalidInput
	}
	now := time.Now()
	product := &entities.Product{
		Name:              strings.TrimSpace(input.Name),
		Description:       strings.TrimSpace(input.Description),
		Category:          strings.TrimSpace(input.Category),
		Price:             input.Price,
		StockQuantity:     input.StockQuantity,
		LowStockThreshold: input.LowStockThreshold,
		Version:           1,
		TimeStamp: entities.TimeStamp{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	if err := u.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}

func (u *ProductUsecase) Update(ctx context.Context, id primitive.ObjectID, input entities.Product) (*entities.Product, error) {
	existing, err := u.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	existing.Name = strings.TrimSpace(input.Name)
	existing.Description = strings.TrimSpace(input.Description)
	existing.Category = strings.TrimSpace(input.Category)
	existing.Price = input.Price
	existing.StockQuantity = input.StockQuantity
	existing.LowStockThreshold = input.LowStockThreshold
	existing.UpdatedAt = time.Now()

	if existing.Name == "" || existing.Price <= 0 || existing.StockQuantity < 0 {
		return nil, entities.ErrInvalidInput
	}

	if err = u.productRepo.Update(ctx, existing); err != nil {
		return nil, err
	}
	existing.Version++
	return existing, nil
}

func (u *ProductUsecase) Delete(ctx context.Context, id primitive.ObjectID) error {
	return u.productRepo.Delete(ctx, id)
}

func (u *ProductUsecase) GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Product, error) {
	return u.productRepo.GetByID(ctx, id)
}

func (u *ProductUsecase) List(ctx context.Context, filter repositories.ProductFilter) ([]entities.Product, int64, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 || filter.PageSize > 100 {
		filter.PageSize = 10
	}
	if filter.SortDir != -1 {
		filter.SortDir = 1
	}
	return u.productRepo.List(ctx, filter)
}
