package repositories

import (
	"context"

	"inventory-management-system/backend/internal/domain/entities"
)

type StockLogRepository interface {
	Create(ctx context.Context, log *entities.StockLog) error
}
