package dto

type ListHistoryPointRequest struct {
	Limit      int32  `json:"limit,omitempty" query:"limit,omitempty"`
	Offset     int32  `json:"offset,omitempty" query:"offset,omitempty"`
	CustomerID string `json:"customer_id,omitempty"`
}
