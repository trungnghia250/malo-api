package dto

import (
	"time"
)

type GetCustomerByIDRequest struct {
	CustomerID string `json:"customer_id" query:"customer_id"`
}

type GetCustomerByIDResponse struct {
	CustomerID   string `json:"customer_id"`
	CustomerName string `json:"customer_name"`
	Gender       string `json:"gender"`
	PhoneNumber  string `json:"phone_number"`
}

type Customer struct {
	CustomerID     string    `bson:"customer_id" json:"customer_id"`
	CustomerName   string    `bson:"customer_name,omitempty" json:"customer_name,omitempty"`
	Gender         string    `bson:"gender,omitempty" json:"gender,omitempty"`
	PhoneNumber    string    `bson:"phone_number,omitempty" json:"phone_number,omitempty"`
	Email          string    `bson:"email,omitempty" json:"email,omitempty"`
	Address        []string  `json:"address,omitempty" bson:"address,omitempty"`
	DateOfBirth    time.Time `json:"date_of_birth,omitempty" bson:"date_of_birth,omitempty"`
	Company        Company   `json:"company,omitempty" bson:"company,omitempty"`
	CustomerType   string    `json:"customer_type,omitempty" bson:"customer_type,omitempty"`
	CustomerSource string    `json:"customer_source,omitempty" bson:"customer_source,omitempty"`
	Tags           []string  `json:"tags,omitempty" bson:"tags,omitempty"`
	AssignFor      []string  `json:"assign_for,omitempty" bson:"assign_for,omitempty"`
	Status         string    `json:"status,omitempty" bson:"status,omitempty"`
	Note           string    `json:"note,omitempty" bson:"note,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ModifiedAt     time.Time `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
	ModifiedBy     string    `json:"modified_by,omitempty" bson:"modified_by,omitempty"`
}

type Company struct {
	CompanyName             string `json:"company_name,omitempty" bson:"company_name,omitempty"`
	Email                   string `json:"email,omitempty" bson:"email,omitempty"`
	PhoneNumber             string `bson:"phone_number,omitempty" json:"phone_number,omitempty"`
	Fax                     string `json:"fax,omitempty" bson:"fax,omitempty"`
	Website                 string `json:"website,omitempty" bson:"website,omitempty"`
	Address                 string `json:"address,omitempty" bson:"address,omitempty"`
	TaxIdentificationNumber string `json:"tax_identification_number,omitempty" bson:"tax_identification_number,omitempty"`
	BusinessLine            string `json:"business_line,omitempty" bson:"business_line,omitempty"`
	RoleOfCustomer          string `json:"role_of_customer,omitempty" bson:"role_of_customer,omitempty"`
	Note                    string `json:"note,omitempty" bson:"note,omitempty"`
}

type ListCustomerRequest struct {
	Limit  int32 `json:"limit,omitempty"`
	Offset int32 `json:"offset,omitempty"`
}
