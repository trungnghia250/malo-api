package dto

import (
	"github.com/trungnghia250/malo-api/service/model"
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
	Address        string    `json:"address,omitempty" bson:"address,omitempty"`
	DateOfBirth    time.Time `json:"date_of_birth,omitempty" bson:"date_of_birth,omitempty"`
	CustomerType   string    `json:"customer_type,omitempty" bson:"customer_type,omitempty"`
	CustomerSource string    `json:"customer_source,omitempty" bson:"customer_source,omitempty"`
	Tags           []string  `json:"tags,omitempty" bson:"tags,omitempty"`
	Note           string    `json:"note,omitempty" bson:"note,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ModifiedAt     time.Time `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
	ModifiedBy     string    `json:"modified_by,omitempty" bson:"modified_by,omitempty"`
}

type ListCustomerRequest struct {
	Limit        int32    `json:"limit,omitempty"`
	Offset       int32    `json:"offset,omitempty"`
	CustomerName []string `json:"customer_name,omitempty" query:"customer_name,omitempty"`
	Phone        string   `json:"phone,omitempty"`
	Address      string   `json:"address,omitempty"`
	CustomerType []string `json:"customer_type,omitempty" query:"customer_type,omitempty"`
	Tags         []string `json:"tags,omitempty" query:"tags,omitempty"`
	Gender       string   `json:"gender,omitempty"`
	Email        string   `json:"email,omitempty"`
	Source       []string `json:"source,omitempty" query:"source,omitempty"`
}

type ListCustomerResponse struct {
	Count int32            `json:"count"`
	Data  []model.Customer `json:"data"`
}

type DeleteCustomersRequest struct {
	CustomerIDs []string `json:"customer_id" query:"customer_ids"`
}

type UpdateListCustomerRequest struct {
	CustomerIDs []string `json:"customer_ids" query:"customer_ids"`
	Tags        []string `json:"tags" query:"tags"`
}
