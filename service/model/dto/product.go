package dto

type GetProductByIDRequest struct {
	ProductID string `json:"product_id" query:"product_id"`
}

type ListProductRequest struct {
	Limit  int32 `json:"limit,omitempty"`
	Offset int32 `json:"offset,omitempty"`
}
