package model

import "time"

type Voucher struct {
	Code             string    `json:"code,omitempty" bson:"code,omitempty"`
	GroupIDs         []string  `json:"group_ids,omitempty" bson:"group_ids,omitempty"`
	DiscountAmount   int32     `json:"discount_amount,omitempty" bson:"discount_amount,omitempty"`
	MinOrderAmount   int32     `json:"min_order_amount,omitempty" bson:"min_order_amount,omitempty"`
	StartAt          time.Time `json:"start_at,omitempty" bson:"start_at,omitempty"`
	ExpireAt         time.Time `json:"expire_at,omitempty" bson:"expire_at,omitempty"`
	Limit            int32     `json:"limit,omitempty" bson:"limit,omitempty"`
	LimitPerCustomer int32     `json:"limit_per_customer,omitempty" bson:"limit_per_customer,omitempty"`
	Note             string    `json:"note,omitempty" bson:"note,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ModifiedAt       time.Time `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
	ModifiedBy       string    `json:"modified_by,omitempty" bson:"modified_by,omitempty"`
	CreatedBy        string    `json:"created_by,omitempty" bson:"created_by,omitempty"`
	TotalCount       int32     `json:"totalCount,omitempty" bson:"totalCount,omitempty"`
}
