package model

import "time"

type SapoOrder struct {
	ID                  int32           `json:"id"`
	ShippingAddress     ShippingAddress `json:"shipping_address"`
	Email               string          `json:"email"`
	DiscountCodes       []DiscountCode  `json:"discount_codes,omitempty"`
	TotalLineItemsPrice float64         `json:"total_line_items_price"`
	ShippingLines       []ShippingLines `json:"shipping_lines"`
	TotalDiscounts      float64         `json:"total_discounts"`
	TotalPrice          float64         `json:"total_price"`
	Status              string          `json:"status"`
	Note                string          `json:"note"`
	CreatedOn           time.Time       `json:"created_on"`
	LineItems           []LineItems     `json:"line_items"`
}

type ShippingAddress struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Address1  string `json:"address1"`
}

type DiscountCode struct {
	Code string `json:"code"`
}

type ShippingLines struct {
	Price float64 `json:"price"`
}

type LineItems struct {
	SKU           string  `json:"sku"`
	Name          string  `json:"name"`
	Quantity      int32   `json:"quantity"`
	Price         float64 `json:"price"`
	TotalDiscount float64 `json:"total_discount"`
}
