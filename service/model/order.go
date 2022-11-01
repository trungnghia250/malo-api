package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Order struct {
	ID                   primitive.ObjectID `json:"_id" bson:"_id"`
	OrderID              string             `json:"order_id" bson:"order_id"`
	CustomerName         string             `json:"customer_name" bson:"customer_name"`
	PhoneNumber          string             `json:"phone_number" bson:"phone_number"`
	Email                string             `json:"email" bson:"email"`
	Address              string             `json:"address" bson:"address"`
	Source               string             `json:"source" bson:"source"`
	Status               string             `json:"status" bson:"status"`
	Items                []Item             `json:"line_items" bson:"line_items"`
	VoucherCode          string             `json:"voucher_code" bson:"voucher_code"`
	TotalLineItemsAmount int32              `json:"total_line_items_amount" bson:"total_line_items_amount"`
	ShippingPrice        int32              `json:"shipping_price" bson:"shipping_price"`
	TotalDiscount        int32              `json:"total_discount" bson:"total_discount"`
	TotalTaxAmount       int32              `json:"total_tax_amount" bson:"total_tax_amount"`
	TotalOrderAmount     int32              `json:"total_order_amount" bson:"total_order_amount"`
	Note                 string             `json:"note" bson:"note"`
	CreateAt             time.Time          `json:"create_at" bson:"create_at"`
	ModifiedAt           time.Time          `json:"modified_at" bson:"modified_at"`
	ModifiedBy           string             `json:"modified_by" bson:"modified_by"`
	TotalCount           int32              `json:"totalCount"`
}

type Item struct {
	ProductName    string `json:"product_name" bson:"product_name"`
	SKU            string `json:"sku" bson:"sku"`
	Quantity       int32  `json:"quantity" bson:"quantity"`
	UnitPrice      int32  `json:"unit_price" bson:"unit_price"`
	TotalDiscount  int32  `json:"total_discount" bson:"total_discount"`
	TotalTaxAmount int32  `json:"total_tax_amount" bson:"total_tax_amount"`
	SubTotal       int32  `json:"sub_total" bson:"sub_total"`
}
