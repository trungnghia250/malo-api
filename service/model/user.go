package model

import "time"

type User struct {
	UserID       string    `bson:"user_id" json:"user_id"`
	Name         string    `bson:"name" json:"name"`
	Email        string    `bson:"email" json:"email"`
	Password     string    `bson:"password" json:"password"`
	Role         string    `bson:"role" json:"role"`
	Token        string    `json:"token" bson:"token"`
	RefreshToken string    `json:"refresh_token" bson:"refresh_token"`
	CreatedAt    time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ModifiedAt   time.Time `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
	ModifiedBy   string    `json:"modified_by,omitempty" bson:"modified_by,omitempty"`
	TotalCount   int32     `json:"totalCount"`
}
