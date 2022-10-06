package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Customer struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	CustomerID   string             `bson:"customer_id" json:"customer_id"`
	CustomerName string             `bson:"customer_name" json:"customer_name"`
	Gender       string             `bson:"gender" bson:"gender"`
	PhoneNumber  string             `bson:"phone_number" bson:"phone_number"`
}
