package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type HistoryPoint struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CustomerID    string             `json:"customer_id,omitempty" bson:"customer_id,omitempty"`
	CustomerName  string             `json:"customer_name,omitempty" bson:"customer_name,omitempty"`
	CustomerPhone string             `json:"customer_phone,omitempty" bson:"customer_phone,omitempty"`
	RewardPoint   int32              `json:"reward_point,omitempty" bson:"reward_point,omitempty"`
	Type          string             `json:"type,omitempty" bson:"type,omitempty"`
	OrderID       string             `json:"order_id,omitempty" bson:"order_id,omitempty"`
	GiftID        string             `json:"gift_id,omitempty" bson:"gift_id,omitempty"`
	Content       string             `json:"content,omitempty" bson:"content,omitempty"`
	CreatedAt     time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	RedeemID      string             `json:"redeem_id,omitempty" bson:"redeem_id,omitempty"`
}
