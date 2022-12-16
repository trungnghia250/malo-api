package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Customer struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CustomerID     string             `bson:"customer_id,omitempty" json:"customer_id,omitempty"`
	CustomerName   string             `bson:"customer_name,omitempty" json:"customer_name,omitempty"`
	Gender         string             `bson:"gender,omitempty" json:"gender,omitempty"`
	PhoneNumber    string             `bson:"phone_number,omitempty" json:"phone_number,omitempty"`
	Email          string             `bson:"email,omitempty" json:"email,omitempty"`
	Address        string             `json:"address,omitempty" bson:"address,omitempty"`
	Province       string             `json:"province,omitempty" bson:"province,omitempty"`
	DateOfBirth    *time.Time         `json:"date_of_birth,omitempty" bson:"date_of_birth,omitempty"`
	CustomerType   string             `json:"customer_type,omitempty" bson:"customer_type,omitempty"`
	CustomerSource string             `json:"customer_source,omitempty" bson:"customer_source,omitempty"`
	Tags           []string           `json:"tags,omitempty" bson:"tags,omitempty"`
	Note           string             `json:"note,omitempty" bson:"note,omitempty"`
	CreatedAt      time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ModifiedAt     time.Time          `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
	ModifiedBy     string             `json:"modified_by,omitempty" bson:"modified_by,omitempty"`
	TotalCount     int32              `json:"totalCount,omitempty"`
	RewardPoint    int32              `json:"reward_point,omitempty" bson:"reward_point,omitempty"`
	RankPoint      int32              `json:"rank_point,omitempty" bson:"rank_point,omitempty"`
	IsNew          bool
}
