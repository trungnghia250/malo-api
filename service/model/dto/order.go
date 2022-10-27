package dto

import (
	"github.com/trungnghia250/malo-api/service/model"
	"time"
)

type GetOrderByIDRequest struct {
	OrderID string `json:"order_id" query:"order_id"`
}

type Order struct {
	OrderID              string    `json:"order_id" bson:"order_id"`
	CustomerName         string    `json:"customer_name,omitempty" bson:"customer_name,omitempty"`
	PhoneNumber          string    `json:"phone_number,omitempty" bson:"phone_number,omitempty"`
	Email                string    `json:"email,omitempty" bson:"email,omitempty"`
	Address              string    `json:"address,omitempty" bson:"address,omitempty"`
	Source               string    `json:"source,omitempty" bson:"source,omitempty"`
	Status               string    `json:"status,omitempty" bson:"status,omitempty"`
	Items                []Item    `json:"line_items,omitempty" bson:"line_items,omitempty"`
	VoucherCode          string    `json:"voucher_code,omitempty" bson:"voucher_code,omitempty"`
	TotalLineItemsAmount int32     `json:"total_line_items_amount,omitempty" bson:"total_line_items_amount,omitempty"`
	ShippingPrice        int32     `json:"shipping_price,omitempty" bson:"shipping_price,omitempty"`
	TotalDiscount        int32     `json:"total_discount,omitempty" bson:"total_discount,omitempty"`
	TotalTaxAmount       int32     `json:"total_tax_amount,omitempty" bson:"total_tax_amount,omitempty"`
	TotalOrderAmount     int32     `json:"total_order_amount,omitempty" bson:"total_order_amount,omitempty"`
	Note                 string    `json:"note,omitempty" bson:"note,omitempty"`
	CreateAt             time.Time `json:"create_at,omitempty" bson:"create_at,omitempty"`
	ModifiedAt           time.Time `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
	ModifiedBy           string    `json:"modified_by,omitempty" bson:"modified_by,omitempty"`
}

type Item struct {
	ProductName    string `json:"product_name,omitempty" bson:"product_name,omitempty"`
	SKU            string `json:"sku,omitempty" bson:"sku,omitempty"`
	Quantity       int32  `json:"quantity,omitempty" bson:"quantity,omitempty"`
	UnitPrice      int32  `json:"unit_price,omitempty" bson:"unit_price,omitempty"`
	TotalDiscount  int32  `json:"total_discount,omitempty" bson:"total_discount,omitempty"`
	TotalTaxAmount int32  `json:"total_tax_amount,omitempty" bson:"total_tax_amount,omitempty"`
	SubTotal       int32  `json:"sub_total,omitempty" bson:"sub_total,omitempty"`
}

type ListOrderResponse struct {
	Count int32         `json:"count"`
	Data  []model.Order `json:"data"`
}
