package dto

import "github.com/trungnghia250/malo-api/service/model"

type ListGiftRequest struct {
	Limit         int32    `json:"limit,omitempty" query:"limit,omitempty"`
	Offset        int32    `json:"offset,omitempty" query:"offset,omitempty"`
	Status        string   `json:"status,omitempty" query:"status,omitempty"`
	CreatedAt     []int32  `json:"created_at,omitempty" query:"created_at,omitempty"`
	Name          []string `json:"name,omitempty" query:"name,omitempty"`
	SKU           []string `json:"sku,omitempty" query:"sku,omitempty"`
	Price         []int32  `json:"price,omitempty" query:"price,omitempty"`
	RewardPoint   []int32  `json:"reward_point,omitempty" query:"reward_point,omitempty"`
	StockAmount   []int32  `json:"stock_amount,omitempty" query:"stock_amount,omitempty"`
	UsedAmount    []int32  `json:"used_amount,omitempty" query:"used_amount,omitempty"`
	ReleaseAmount []int32  `json:"release_amount,omitempty" query:"release_amount,omitempty"`
	Category      []string `json:"category,omitempty" query:"category,omitempty"`
	GiftIDs       []string `json:"gift_ids,omitempty" query:"gift_ids,omitempty"`
}

type ListGiftResponse struct {
	Count int32        `json:"count"`
	Data  []model.Gift `json:"data"`
}

type GetGiftByIDRequest struct {
	GiftID string `json:"gift_id" query:"gift_id"`
}

type DeleteGiftsRequest struct {
	IDs []string `json:"ids" query:"ids"`
}
