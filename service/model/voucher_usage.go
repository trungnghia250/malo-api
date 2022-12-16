package model

import (
	"time"
)

type VoucherUsage struct {
	ID             string    `json:"_id,omitempty" bson:"_id,omitempty"`
	Code           string    `json:"code,omitempty" bson:"code,omitempty"`
	CustomerName   string    `json:"customer_name,omitempty" bson:"customer_name,omitempty"`
	Phone          string    `json:"phone,omitempty" bson:"phone,omitempty"`
	OrderID        string    `json:"order_id,omitempty" bson:"order_id,omitempty"`
	DiscountAmount int32     `json:"discount_amount,omitempty" bson:"discount_amount,omitempty"`
	ProofURL       string    `json:"proof_url,omitempty" bson:"proof_url,omitempty"`
	Note           string    `json:"note,omitempty" bson:"note,omitempty"`
	UsedDate       int32     `json:"used_date,omitempty" bson:"used_date,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ModifiedAt     time.Time `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
	ModifiedBy     string    `json:"modified_by,omitempty" bson:"modified_by,omitempty"`
	CreatedBy      string    `json:"created_by,omitempty" bson:"created_by,omitempty"`
	TotalCount     int32     `json:"totalCount,omitempty" bson:"totalCount,omitempty"`
	TotalDiscount  int32     `json:"total_discount,omitempty" bson:"total_discount,omitempty"`
}
