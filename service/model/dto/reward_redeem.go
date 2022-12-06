package dto

import "github.com/trungnghia250/malo-api/service/model"

type ListRewardRedeemRequest struct {
	Limit  int32    `json:"limit,omitempty" query:"limit,omitempty"`
	Offset int32    `json:"offset,omitempty" query:"offset,omitempty"`
	IDs    []string `json:"ids,omitempty"`
}

type ListRewardRedeemResponse struct {
	Count int32                `json:"count"`
	Data  []model.RewardRedeem `json:"data"`
}

type GetRewardRedeemByIDRequest struct {
	ID string `json:"id" query:"id"`
}

type DeleteRedeemsRequest struct {
	IDs []string `json:"ids" query:"ids"`
}
