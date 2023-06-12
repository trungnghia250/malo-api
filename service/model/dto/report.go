package dto

import "time"

type CustomerReport struct {
	Name             string    `json:"name,omitempty" bson:"name,omitempty"`
	Phone            string    `json:"_id,omitempty" bson:"_id,omitempty"`
	Email            string    `json:"email,omitempty" bson:"email,omitempty"`
	TotalOrders      int32     `json:"total_orders,omitempty" bson:"total_orders,omitempty"`
	SuccessOrders    int32     `json:"success_orders,omitempty" bson:"success_orders,omitempty"`
	ProcessingOrders int32     `json:"processing_orders,omitempty" bson:"processing_orders,omitempty"`
	CancelOrders     int32     `json:"cancel_orders,omitempty" bson:"cancel_orders,omitempty"`
	TotalRevenue     int32     `json:"total_revenue,omitempty" bson:"total_revenue,omitempty"`
	New              int32     `json:"new" bson:"new,omitempty"`
	Return           int32     `json:"return" bson:"return"`
	Date             time.Time `json:"date,omitempty" bson:"date,omitempty"`
}

type CustomerReportResponse struct {
	Data  []CustomerReport `json:"data"`
	Count int32            `json:"count"`
	Total CustomerReport   `json:"total"`
}

type GetReportRequest struct {
	Type            string   `json:"type" query:"type"`
	StartTime       string   `json:"start_time" query:"start_time"`
	EndTime         string   `json:"end_time" query:"end_time"`
	Limit           int32    `json:"limit,omitempty" query:"limit,omitempty"`
	Offset          int32    `json:"offset,omitempty" query:"offset,omitempty"`
	Name            []string `json:"name,omitempty" query:"name,omitempty"`
	SKU             []string `json:"sku,omitempty" query:"sku,omitempty"`
	Phone           []string `json:"phone,omitempty" query:"phone,omitempty"`
	Email           []string `json:"email,omitempty" query:"email,omitempty"`
	TotalOrders     []int32  `json:"total_orders,omitempty" query:"total_orders,omitempty"`
	TotalSales      []int32  `json:"total_sales,omitempty" query:"total_sales,omitempty"`
	TotalSuccess    []int32  `json:"total_success,omitempty" query:"total_success,omitempty"`
	TotalProcessing []int32  `json:"total_processing,omitempty" query:"total_processing,omitempty"`
	TotalCancel     []int32  `json:"total_cancel,omitempty" query:"total_cancel,omitempty"`
	TotalRevenue    []int32  `json:"total_revenue,omitempty" query:"total_revenue,omitempty"`
	Export          bool     `json:"export" query:"export"`
	GroupID         string   `json:"group_id,omitempty" query:"group_id,omitempty"`
}

type ProductReport struct {
	Name             string `json:"name,omitempty" bson:"name,omitempty"`
	SKU              string `json:"_id,omitempty" bson:"_id,omitempty"`
	TotalOrders      int32  `json:"total_orders,omitempty" bson:"total_orders,omitempty"`
	TotalSales       int32  `json:"total_sales,omitempty" bson:"total_sales,omitempty"`
	SuccessOrders    int32  `json:"success_orders,omitempty" bson:"success_orders,omitempty"`
	ProcessingOrders int32  `json:"processing_orders,omitempty" bson:"processing_orders,omitempty"`
	CancelOrders     int32  `json:"cancel_orders,omitempty" bson:"cancel_orders,omitempty"`
	TotalRevenue     int32  `json:"total_revenue,omitempty" bson:"total_revenue,omitempty"`
}

type ProductReportResponse struct {
	Data  []ProductReport `json:"data"`
	Count int32           `json:"count"`
	Total ProductReport   `json:"total"`
}

type GetDashBoardRequest struct {
	StartTime  string   `json:"start_time" query:"start_time"`
	EndTime    string   `json:"end_time" query:"end_time"`
	CustomerID []string `json:"customer_id,omitempty" query:"customer_id,omitempty"`
	GroupID    []string `json:"group_id,omitempty" query:"group_id,omitempty"`
	Type       string   `json:"type,omitempty" query:"type,omitempty"`
}
