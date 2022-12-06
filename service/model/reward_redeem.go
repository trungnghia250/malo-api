package model

import "time"

type RewardRedeem struct {
	ID           string    `json:"_id,omitempty" bson:"_id,omitempty"`
	CustomerName string    `json:"customer_name,omitempty" bson:"customer_name,omitempty"`
	Phone        string    `json:"phone,omitempty" bson:"phone,omitempty"`
	RewardPoint  int32     `json:"reward_point,omitempty" bson:"reward_point,omitempty"`
	RewardType   string    `json:"reward_type,omitempty" bson:"reward_type,omitempty"`
	GiftID       string    `json:"gift_id,omitempty" bson:"gift_id,omitempty"`
	OrderID      string    `json:"order_id,omitempty" bson:"order_id,omitempty"`
	RewardValue  int32     `json:"reward_value,omitempty" bson:"reward_value,omitempty"`
	RedeemDate   time.Time `json:"redeem_date,omitempty" bson:"redeem_date,omitempty"`
	Note         string    `json:"note,omitempty" bson:"note,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ModifiedAt   time.Time `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
	ModifiedBy   string    `json:"modified_by,omitempty" bson:"modified_by,omitempty"`
	CreatedBy    string    `json:"created_by,omitempty" bson:"created_by,omitempty"`
	TotalCount   int32     `json:"totalCount,omitempty" bson:"totalCount,omitempty"`
}
