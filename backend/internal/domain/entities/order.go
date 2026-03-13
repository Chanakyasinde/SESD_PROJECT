package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type OrderItem struct {
	ProductID   primitive.ObjectID `bson:"productId" json:"productId"`
	ProductName string             `bson:"productName" json:"productName"`
	Quantity    int64              `bson:"quantity" json:"quantity"`
	UnitPrice   float64            `bson:"unitPrice" json:"unitPrice"`
	SubTotal    float64            `bson:"subTotal" json:"subTotal"`
}

type Order struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CustomerID  primitive.ObjectID `bson:"customerId" json:"customerId"`
	CreatedBy   primitive.ObjectID `bson:"createdBy" json:"createdBy"`
	Status      string             `bson:"status" json:"status"`
	Items       []OrderItem        `bson:"items" json:"items"`
	TotalAmount float64            `bson:"totalAmount" json:"totalAmount"`
	TimeStamp   `bson:",inline"`
}
