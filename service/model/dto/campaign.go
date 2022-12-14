package dto

import "github.com/trungnghia250/malo-api/service/model"

type ListCampaignRequest struct {
	Limit     int32    `json:"limit,omitempty"`
	Offset    int32    `json:"offset,omitempty"`
	Status    []string `json:"status,omitempty"`
	CreatedAt []int32  `json:"created_at,omitempty" query:"created_at,omitempty"`
	SendAt    []int32  `json:"send_at,omitempty" query:"send_at,omitempty"`
}

type GetCampaignByIDRequest struct {
	CampaignID string `json:"campaign_id" query:"campaign_id"`
}

type DeleteCampaignsRequest struct {
	CampaignIDs []string `json:"campaign_ids" query:"campaign_ids"`
}

type ListCampaignResponse struct {
	Count int32            `json:"count"`
	Data  []model.Campaign `json:"data"`
}
