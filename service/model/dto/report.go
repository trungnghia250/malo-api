package dto

type CustomerReport struct {
	Name             string `json:"name,omitempty" bson:"name,omitempty"`
	Phone            string `json:"_id,omitempty" bson:"_id,omitempty"`
	Email            string `json:"email,omitempty" bson:"email,omitempty"`
	TotalOrders      int32  `json:"total_orders,omitempty" bson:"total_orders,omitempty"`
	SuccessOrders    int32  `json:"success_orders,omitempty" bson:"success_orders,omitempty"`
	ProcessingOrders int32  `json:"processing_orders,omitempty" bson:"processing_orders,omitempty"`
	CancelOrders     int32  `json:"cancel_orders,omitempty" bson:"cancel_orders,omitempty"`
	TotalRevenue     int32  `json:"total_revenue,omitempty" bson:"total_revenue,omitempty"`
}

type CustomerReportResponse struct {
	Data  []CustomerReport `json:"data"`
	Count int32            `json:"count"`
	Total CustomerReport   `json:"total"`
}

type GetReportRequest struct {
	Type      string `json:"type" query:"type"`
	StartTime string `json:"start_time" query:"start_time"`
	EndTime   string `json:"end_time" query:"end_time"`
	Limit     int32  `json:"limit,omitempty" query:"limit,omitempty"`
	Offset    int32  `json:"offset,omitempty" query:"offset,omitempty"`
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
