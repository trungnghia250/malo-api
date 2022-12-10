package dto

type GetPartnerConfigByIDRequest struct {
	PartnerCode string `json:"partner_code" query:"partner_code"`
}

type UploadResponse struct {
	URL string `json:"url"`
}
