package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name              string             `bson:"name" json:"name"`
	Description       string             `bson:"description" json:"description"`
	Category          string             `bson:"category" json:"category"`
	Price             float64            `bson:"price" json:"price"`
	StockQuantity     int64              `bson:"stockQuantity" json:"stockQuantity"`
	LowStockThreshold int64              `bson:"lowStockThreshold" json:"lowStockThreshold"`
	Version           int64              `bson:"version" json:"version"`
	TimeStamp         `bson:",inline"`
}
