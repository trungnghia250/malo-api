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
	Company        Company            `json:"company,omitempty" bson:"company,omitempty"`
	CustomerType   string             `json:"customer_type" bson:"customer_type"`
	CustomerSource string             `json:"customer_source" bson:"customer_source"`
	Tags           []string           `json:"tags" bson:"tags"`
	AssignFor      []string           `json:"assign_for" bson:"assign_for"`
	Status         string             `json:"status" bson:"status"`
	Note           string             `json:"note" bson:"note"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	ModifiedAt     time.Time          `json:"modified_at" bson:"modified_at"`
	ModifiedBy     string             `json:"modified_by" bson:"modified_by"`
}

type Company struct {
	CompanyName             string `json:"company_name" bson:"company_name"`
	Email                   string `json:"email" bson:"email"`
	PhoneNumber             string `bson:"phone_number" json:"phone_number"`
	Fax                     string `json:"fax" bson:"fax"`
	Website                 string `json:"website" bson:"website"`
	Address                 string `json:"address" bson:"address"`
	TaxIdentificationNumber string `json:"tax_identification_number" bson:"tax_identification_number"`
	BusinessLine            string `json:"business_line" bson:"business_line"`
	RoleOfCustomer          string `json:"role_of_customer" bson:"role_of_customer"`
	Note                    string `json:"note" bson:"note"`
}
