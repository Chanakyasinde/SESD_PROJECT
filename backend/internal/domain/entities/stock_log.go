package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type StockLog struct {
	ID              primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	ProductID       primitive.ObjectID  `bson:"productId" json:"productId"`
	OrderID         *primitive.ObjectID `bson:"orderId,omitempty" json:"orderId,omitempty"`
	ChangedBy       primitive.ObjectID  `bson:"changedBy" json:"changedBy"`
	QuantityChanged int64               `bson:"quantityChanged" json:"quantityChanged"`
	ChangeType      string              `bson:"changeType" json:"changeType"`
	Reason          string              `bson:"reason" json:"reason"`
	TimeStamp       `bson:",inline"`
}
