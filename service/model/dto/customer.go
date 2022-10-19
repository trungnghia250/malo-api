package dto

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
	CustomerID   string `bson:"customer_id,omitempty" json:"customer_id,omitempty"`
	CustomerName string `bson:"customer_name,omitempty" json:"customer_name,omitempty"`
	Gender       string `bson:"gender,omitempty" bson:"gender,omitempty"`
	PhoneNumber  string `bson:"phone_number,omitempty" bson:"phone_number,omitempty"`
}

type ListCustomerRequest struct {
	Limit  int32 `json:"limit,omitempty"`
	Offset int32 `json:"offset,omitempty"`
}
