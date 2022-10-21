package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Product struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	ProductID   string             `json:"product_id" bson:"product_id"`
	ProductName string             `json:"product_name" bson:"product_name"`
	SKU         string             `json:"sku" bson:"sku"`
	Description string             `json:"description" bson:"description"`
	Image       string             `json:"image" bson:"image"`
	Category    string             `json:"category" bson:"category"`
	Note        string             `json:"note" bson:"note"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	ModifiedAt  time.Time          `json:"modified_at" bson:"modified_at"`
	ModifiedBy  string             `json:"modified_by" bson:"modified_by"`
}
