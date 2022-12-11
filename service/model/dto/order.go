package dto

import (
	"github.com/trungnghia250/malo-api/service/model"
	"mime/multipart"
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
	Province             string    `json:"province,omitempty" bson:"province,omitempty"`
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
	Gender               string    `json:"gender,omitempty" bson:"gender,omitempty"`
	CreatedAt            time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
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
	Subtotal       int32  `json:"subtotal,omitempty" bson:"subtotal,omitempty"`
}

type ListOrderResponse struct {
	Count int32         `json:"count"`
	Data  []model.Order `json:"data"`
}

type ListOrderRequest struct {
	Limit                int32    `json:"limit,omitempty"`
	Offset               int32    `json:"offset,omitempty"`
	CustomerName         []string `json:"customer_name,omitempty" query:"customer_name,omitempty"`
	Email                string   `json:"email,omitempty"`
	Source               []string `json:"source,omitempty"`
	TotalLineItemsAmount []int32  `json:"total_line_items_amount,omitempty" query:"total_line_items_amount,omitempty"`
	TotalDiscount        []int32  `json:"total_discount,omitempty" query:"total_discount,omitempty"`
	TotalOrderAmount     []int32  `json:"total_order_amount,omitempty" query:"total_order_amount,omitempty"`
	Phone                string   `json:"phone,omitempty"`
	Address              string   `json:"address,omitempty"`
	VoucherCode          string   `json:"voucher_code,omitempty" query:"voucher_code,omitempty"`
	ShippingPrice        []int32  `json:"shipping_price,omitempty" query:"shipping_price,omitempty"`
	TotalTax             []int32  `json:"total_tax,omitempty" query:"total_tax,omitempty"`
	Status               []string `json:"status,omitempty" query:"status,omitempty"`
	OrderIDs             []string `json:"order_ids,omitempty" query:"order_ids,omitempty"`
}

type ImportOrderRequest struct {
	CheckDupCol []string              `form:"check_dup_col"`
	Action      string                `form:"action"`
	File        *multipart.FileHeader `form:"file"`
}

type ImportOrderResponse struct {
	Scan    int32         `json:"scan"`
	Success int32         `json:"success"`
	Insert  int32         `json:"insert"`
	Update  int32         `json:"update"`
	Ignore  int32         `json:"ignore"`
	Data    []model.Order `json:"data"`
}

type DeleteOrdersRequest struct {
	OrderIDs []string `json:"order_ids" query:"order_ids"`
}

type ExportOrderRequest struct {
	OrderIDs []string `json:"order_ids,omitempty" query:"order_ids,omitempty"`
	Filter   string   `json:"filter,omitempty" query:"filter,omitempty"`
}

type ExportOrderResponse struct {
	FileOrderURL string `json:"file_order_url"`
}

type SyncOrderRequest struct {
	Source      string `json:"source,omitempty"`
	StartTime   string `query:"start_time,omitempty"`
	EndTime     string `query:"end_time,omitempty"`
	CheckDupCol string `json:"check_dup_col"`
	ArrayCol    []string
	Action      string `json:"action"`
}

type SapoOrdersResponse struct {
	Orders []model.SapoOrder `json:"orders"`
}
