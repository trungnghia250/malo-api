package dto

import "github.com/trungnghia250/malo-api/service/model"

type GetProductByIDRequest struct {
	ProductID string `json:"product_id" query:"product_id"`
}

type ListProductRequest struct {
	Limit  int32 `json:"limit,omitempty"`
	Offset int32 `json:"offset,omitempty"`
}

type ListProductResponse struct {
	Count int32           `json:"count"`
	Data  []model.Product `json:"data"`
}
