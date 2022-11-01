package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Customer struct {
	ID             primitive.ObjectID `json:"_id" bson:"_id"`
	CustomerID     string             `bson:"customer_id" json:"customer_id"`
	CustomerName   string             `bson:"customer_name" json:"customer_name"`
	Gender         string             `bson:"gender" json:"gender"`
	PhoneNumber    string             `bson:"phone_number" json:"phone_number"`
	Email          string             `bson:"email" json:"email"`
	Address        string             `json:"address" bson:"address"`
	DateOfBirth    time.Time          `json:"date_of_birth" bson:"date_of_birth"`
	CustomerType   string             `json:"customer_type" bson:"customer_type"`
	CustomerSource string             `json:"customer_source" bson:"customer_source"`
	Tags           []string           `json:"tags" bson:"tags"`
	Note           string             `json:"note" bson:"note"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	ModifiedAt     time.Time          `json:"modified_at" bson:"modified_at"`
	ModifiedBy     string             `json:"modified_by" bson:"modified_by"`
	TotalCount     int32              `json:"totalCount"`
}
