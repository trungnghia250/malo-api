package model

import "time"

type CustomerReport struct {
	Phone         string    `json:"phone,omitempty" bson:"phone,omitempty"`
	Date          time.Time `json:"date,omitempty" bson:"date,omitempty"`
	Name          string    `json:"name,omitempty" bson:"name,omitempty"`
	Email         string    `json:"email,omitempty" bson:"email,omitempty"`
	TotalOrders   int32     `json:"total_orders,omitempty" bson:"total_orders,omitempty"`
	SuccessOrders int32     `json:"success_orders,omitempty" bson:"success_orders,omitempty"`
	ProcessOrders int32     `json:"processing_orders,omitempty" bson:"processing_orders,omitempty"`
	CancelOrders  int32     `json:"cancel_orders,omitempty" bson:"cancel_orders,omitempty"`
	Revenue       int32     `json:"revenue,omitempty" bson:"revenue,omitempty"`
}

type ProductReport struct {
	Date          time.Time `json:"date,omitempty" bson:"date,omitempty"`
	SKU           string    `json:"sku,omitempty" bson:"sku,omitempty"`
	Name          string    `json:"name,omitempty" bson:"name,omitempty"`
	TotalSales    int32     `json:"total_sales,omitempty" bson:"total_sales,omitempty"`
	TotalOrders   int32     `json:"total_orders,omitempty" bson:"total_orders,omitempty"`
	SuccessOrders int32     `json:"success_orders,omitempty" bson:"success_orders,omitempty"`
	ProcessOrders int32     `json:"processing_orders,omitempty" bson:"processing_orders,omitempty"`
	CancelOrders  int32     `json:"cancel_orders,omitempty" bson:"cancel_orders,omitempty"`
	Revenue       int32     `json:"revenue,omitempty" bson:"revenue,omitempty"`
}
