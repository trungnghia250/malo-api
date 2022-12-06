package model

import "time"

type Voucher struct {
	ID               string    `json:"_id,omitempty" bson:"_id,omitempty"`
	GroupIDs         []string  `json:"group_ids,omitempty" bson:"group_ids,omitempty"`
	DiscountAmount   int32     `json:"discount_amount,omitempty" bson:"discount_amount,omitempty"`
	MinOrderAmount   int32     `json:"min_order_amount,omitempty" bson:"min_order_amount,omitempty"`
	StartAt          int32     `json:"start_at,omitempty" bson:"start_at,omitempty"`
	ExpireAt         int32     `json:"expire_at,omitempty" bson:"expire_at,omitempty"`
	LimitUsage       int32     `json:"limit_usage,omitempty" bson:"limit_usage,omitempty"`
	Status           string    `json:"status,omitempty" bson:"status,omitempty"`
	LimitPerCustomer int32     `json:"limit_per_customer,omitempty" bson:"limit_per_customer,omitempty"`
	Note             string    `json:"note,omitempty" bson:"note,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ModifiedAt       time.Time `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
	ModifiedBy       string    `json:"modified_by,omitempty" bson:"modified_by,omitempty"`
	CreatedBy        string    `json:"created_by,omitempty" bson:"created_by,omitempty"`
	TotalCount       int32     `json:"totalCount,omitempty" bson:"totalCount,omitempty"`
}
