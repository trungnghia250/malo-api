package model

import "time"

type User struct {
	UserID       string    `bson:"user_id,omitempty" json:"user_id,omitempty"`
	Name         string    `bson:"name,omitempty" json:"name,omitempty"`
	Email        string    `bson:"email,omitempty" json:"email,omitempty"`
	Password     string    `bson:"password,omitempty" json:"password,omitempty"`
	Role         string    `bson:"role,omitempty" json:"role,omitempty"`
	Token        string    `json:"token,omitempty" bson:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty" bson:"refresh_token,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ModifiedAt   time.Time `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
	ModifiedBy   string    `json:"modified_by,omitempty" bson:"modified_by,omitempty"`
	TotalCount   int32     `json:"totalCount,omitempty"`
}
