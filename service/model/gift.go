package model

import "time"

type Gift struct {
	ID            string    `json:"_id,omitempty" bson:"_id,omitempty"`
	GroupIDs      []string  `json:"group_ids,omitempty" bson:"group_ids,omitempty"`
	Name          string    `json:"name,omitempty" bson:"name,omitempty"`
	SKU           string    `json:"sku,omitempty" bson:"sku,omitempty"`
	ImageURL      string    `json:"image_url,omitempty" bson:"image_url,omitempty"`
	Price         int32     `json:"price,omitempty" bson:"price,omitempty"`
	RewardPoint   int32     `json:"reward_point,omitempty" bson:"reward_point,omitempty"`
	StockAmount   int32     `json:"stock_amount,omitempty" bson:"stock_amount,omitempty"`
	UsedAmount    int32     `json:"used_amount,omitempty" bson:"used_amount,omitempty"`
	ReleaseAmount int32     `json:"release_amount,omitempty" bson:"release_amount,omitempty"`
	Category      string    `json:"category,omitempty" bson:"category,omitempty"`
	Note          string    `json:"note,omitempty" bson:"note,omitempty"`
	Status        string    `json:"status,omitempty" bson:"status,omitempty"`
	StartAt       int32     `json:"start_at,omitempty" bson:"start_at,omitempty"`
	ExpireAt      int32     `json:"expire_at,omitempty" bson:"expire_at,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ModifiedAt    time.Time `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
	ModifiedBy    string    `json:"modified_by,omitempty" bson:"modified_by,omitempty"`
	CreatedBy     string    `json:"created_by,omitempty" bson:"created_by,omitempty"`
	TotalCount    int32     `json:"totalCount,omitempty" bson:"totalCount,omitempty"`
	GroupNames    []string  `json:"group_names,omitempty"`
}
