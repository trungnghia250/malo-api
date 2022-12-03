package dto

import "github.com/trungnghia250/malo-api/service/model"

type ListTemplateRequest struct {
	Limit  int32 `json:"limit,omitempty"`
	Offset int32 `json:"offset,omitempty"`
}

type ListTemplateResponse struct {
	Count int32            `json:"count"`
	Data  []model.Template `json:"data"`
}

type GetTemplateByIDRequest struct {
	TemplateID string `json:"template_id" query:"template_id"`
}

type DeleteTemplatesRequest struct {
	IDs []string `json:"ids" query:"ids"`
}
