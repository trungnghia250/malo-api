package model

import "time"

type CustomerGroup struct {
	ID          string    `json:"_id,omitempty" bson:"_id,omitempty"`
	GroupName   string    `json:"group_name,omitempty" bson:"group_name,omitempty"`
	CustomerIDs []string  `json:"customer_ids,omitempty" bson:"customer_ids,omitempty"`
	Note        string    `json:"note,omitempty" bson:"note,omitempty"`
	ProductIDs  []string  `json:"product_ids,omitempty" bson:"product_ids,omitempty"`
	CreatedBy   string    `json:"created_by,omitempty" bson:"created_by,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ModifiedAt  time.Time `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
	ModifiedBy  string    `json:"modified_by,omitempty" bson:"modified_by,omitempty"`
	TotalCount  int32     `json:"totalCount,omitempty" bson:"totalCount,omitempty"`
}
