package model

import "time"

type Campaign struct {
	ID               string    `json:"_id,omitempty" bson:"_id,omitempty"`
	Title            string    `json:"title,omitempty" bson:"title,omitempty"`
	Channel          string    `json:"channel,omitempty" bson:"channel,omitempty"`
	Type             string    `json:"type,omitempty" bson:"type,omitempty"`
	Status           string    `json:"status,omitempty" bson:"status,omitempty"`
	SendAt           int       `json:"send_at,omitempty" bson:"send_at,omitempty"`
	Note             string    `json:"note,omitempty" bson:"note,omitempty"`
	VoucherCode      string    `json:"voucher_code,omitempty" bson:"voucher_code,omitempty"`
	CustomerGroupIDs []string  `json:"customer_group_ids,omitempty" bson:"customer_group_ids,omitempty"`
	Message          string    `json:"message,omitempty" bson:"message,omitempty"`
	TotalCount       int32     `json:"totalCount,omitempty" bson:"totalCount,omitempty"`
	CreatedBy        string    `bson:"created_by,omitempty" json:"created_by,omitempty"`
	ModifiedBy       string    `json:"modified_by,omitempty" bson:"modified_by,omitempty"`
	ModifiedAt       time.Time `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}
