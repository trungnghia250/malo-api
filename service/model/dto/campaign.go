package dto

import "github.com/trungnghia250/malo-api/service/model"

type ListCampaignRequest struct {
	Limit  int32 `json:"limit,omitempty"`
	Offset int32 `json:"offset,omitempty"`
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
