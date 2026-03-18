package usecase

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"inventory-management-system/backend/internal/domain/entities"
	"inventory-management-system/backend/internal/domain/repositories"
	"inventory-management-system/backend/internal/domain/usecases"
)

type OrderUsecase struct {
	orderRepo    repositories.OrderRepository
	productRepo  repositories.ProductRepository
	stockLogRepo repositories.StockLogRepository
}

func NewOrderUsecase(orderRepo repositories.OrderRepository, productRepo repositories.ProductRepository, stockLogRepo repositories.StockLogRepository) usecases.OrderUsecase {
	return &OrderUsecase{orderRepo: orderRepo, productRepo: productRepo, stockLogRepo: stockLogRepo}
}

func (u *OrderUsecase) Create(ctx context.Context, actorID primitive.ObjectID, input usecases.CreateOrderInput) (*entities.Order, error) {
	if len(input.Items) == 0 {
		return nil, entities.ErrInvalidInput
	}

	customerID, err := primitive.ObjectIDFromHex(input.CustomerID)
	if err != nil {
		return nil, entities.ErrInvalidInput
	}

	items := make([]entities.OrderItem, 0, len(input.Items))
	total := 0.0
	for _, item := range input.Items {
		if item.Quantity <= 0 {
			return nil, entities.ErrInvalidInput
		}
		productID, err := primitive.ObjectIDFromHex(item.ProductID)
		if err != nil {
			return nil, entities.ErrInvalidInput
		}

		product, err := u.productRepo.GetByID(ctx, productID)
		if err != nil {
			return nil, err
		}
		if product.StockQuantity < item.Quantity {
			return nil, fmt.Errorf("%w: insufficient stock for %s", entities.ErrConflict, product.Name)
		}

		if err = u.productRepo.DecreaseStockWithVersion(ctx, productID, item.Quantity, product.Version); err != nil {
			return nil, err
		}

		subtotal := float64(item.Quantity) * product.Price
		total += subtotal
		items = append(items, entities.OrderItem{
			ProductID:   productID,
			ProductName: product.Name,
			Quantity:    item.Quantity,
			UnitPrice:   product.Price,
			SubTotal:    subtotal,
		})
	}

	now := time.Now()
	order := &entities.Order{
		CustomerID:  customerID,
		CreatedBy:   actorID,
		Status:      entities.OrderConfirmed,
		Items:       items,
		TotalAmount: total,
		TimeStamp: entities.TimeStamp{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	if err := u.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	for _, item := range items {
		orderID := order.ID
		log := &entities.StockLog{
			ProductID:       item.ProductID,
			OrderID:         &orderID,
			ChangedBy:       actorID,
			QuantityChanged: -item.Quantity,
			ChangeType:      entities.StockDeduction,
			Reason:          "order_confirmed",
			TimeStamp:       entities.TimeStamp{CreatedAt: now, UpdatedAt: now},
		}
		if err := u.stockLogRepo.Create(ctx, log); err != nil {
			return nil, err
		}
	}

	return order, nil
}

func (u *OrderUsecase) UpdateStatus(ctx context.Context, actorID primitive.ObjectID, id primitive.ObjectID, status string) error {
	order, err := u.orderRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !isValidTransition(order.Status, status) {
		return entities.ErrInvalidTransition
	}

	if status == entities.OrderCancelled {
		now := time.Now()
		for _, item := range order.Items {
			if err = u.productRepo.IncreaseStock(ctx, item.ProductID, item.Quantity); err != nil {
				return err
			}
			orderID := order.ID
			log := &entities.StockLog{
				ProductID:       item.ProductID,
				OrderID:         &orderID,
				ChangedBy:       actorID,
				QuantityChanged: item.Quantity,
				ChangeType:      entities.StockRestock,
				Reason:          "order_cancelled",
				TimeStamp:       entities.TimeStamp{CreatedAt: now, UpdatedAt: now},
			}
			if err = u.stockLogRepo.Create(ctx, log); err != nil {
				return err
			}
		}
	}

	return u.orderRepo.UpdateStatus(ctx, id, status)
}

func (u *OrderUsecase) GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Order, error) {
	return u.orderRepo.GetByID(ctx, id)
}

func (u *OrderUsecase) List(ctx context.Context, filter repositories.OrderFilter) ([]entities.Order, int64, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 || filter.PageSize > 100 {
		filter.PageSize = 10
	}
	if filter.SortDir != -1 {
		filter.SortDir = 1
	}
	return u.orderRepo.List(ctx, filter)
}

func isValidTransition(from, to string) bool {
	if from == to {
		return true
	}
	allowed := map[string][]string{
		entities.OrderPending:   {entities.OrderConfirmed, entities.OrderCancelled},
		entities.OrderConfirmed: {entities.OrderShipped, entities.OrderCancelled},
		entities.OrderShipped:   {entities.OrderDelivered},
		entities.OrderDelivered: {},
		entities.OrderCancelled: {},
	}
	for _, next := range allowed[from] {
		if next == to {
			return true
		}
	}
	return false
}
