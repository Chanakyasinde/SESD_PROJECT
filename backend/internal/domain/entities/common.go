package entities

import "time"

const (
	RoleAdmin    = "admin"
	RoleStaff    = "staff"
	RoleCustomer = "customer"
)

const (
	OrderPending   = "pending"
	OrderConfirmed = "confirmed"
	OrderShipped   = "shipped"
	OrderDelivered = "delivered"
	OrderCancelled = "cancelled"
)

const (
	StockDeduction  = "deduction"
	StockRestock    = "restock"
	StockAdjustment = "adjustment"
)

type Pagination struct {
	Page     int64
	PageSize int64
	SortBy   string
	SortDir  int
}

type TimeStamp struct {
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
