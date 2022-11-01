package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Product struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ProductID   string             `json:"product_id,omitempty" bson:"product_id,omitempty"`
	ProductName string             `json:"product_name,omitempty" bson:"product_name,omitempty"`
	SKU         string             `json:"sku,omitempty" bson:"sku,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Image       string             `json:"image,omitempty" bson:"image,omitempty"`
	Category    string             `json:"category,omitempty" bson:"category,omitempty"`
	Note        string             `json:"note,omitempty" bson:"note,omitempty"`
	CreatedAt   time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ModifiedAt  time.Time          `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
	ModifiedBy  string             `json:"modified_by,omitempty" bson:"modified_by,omitempty"`
	TotalCount  int32              `json:"totalCount"`
}
