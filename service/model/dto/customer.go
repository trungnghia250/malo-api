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
