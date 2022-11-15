package dto

type GetPartnerConfigByIDRequest struct {
	PartnerCode string `json:"partner_code" query:"partner_code"`
}
