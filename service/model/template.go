package model

import "time"

type Template struct {
	ID         string    `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string    `json:"name,omitempty" bson:"name,omitempty"`
	Message    string    `json:"message,omitempty" bson:"message,omitempty"`
	Type       string    `json:"type,omitempty" bson:"type,omitempty"`
	Note       string    `json:"note,omitempty" bson:"note,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ModifiedAt time.Time `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
	ModifiedBy string    `json:"modified_by,omitempty" bson:"modified_by,omitempty"`
	CreatedBy  string    `json:"created_by,omitempty" bson:"created_by,omitempty"`
	TotalCount int32     `json:"totalCount,omitempty" bson:"totalCount,omitempty"`
}
