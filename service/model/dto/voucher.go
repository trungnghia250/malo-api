package dto

import "github.com/trungnghia250/malo-api/service/model"

type ListVoucherRequest struct {
	Limit            int32    `json:"limit,omitempty" query:"limit,omitempty"`
	Offset           int32    `json:"offset,omitempty" query:"offset,omitempty"`
	Code             []string `json:"code,omitempty" query:"code,omitempty"`
	DiscountAmount   []int32  `json:"discount_amount,omitempty" query:"discount_amount,omitempty"`
	MinOrderAmount   []int32  `json:"min_order_amount,omitempty" query:"min_order_amount,omitempty"`
	StartAt          []int32  `json:"start_at,omitempty" query:"start_at,omitempty"`
	ExpireAt         []int32  `json:"expire_at,omitempty" query:"expire_at,omitempty"`
	Status           string   `json:"status,omitempty" query:"status,omitempty"`
	LimitUsage       []int32  `json:"limit_usage,omitempty" query:"limit_usage,omitempty"`
	LimitPerCustomer []int32  `json:"limit_per_customer,omitempty" query:"limit_per_customer,omitempty"`
	CreatedAt        []int32  `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type ListVoucherResponse struct {
	Count int32           `json:"count"`
	Data  []model.Voucher `json:"data"`
}

type GetVoucherByIDRequest struct {
	ID string `json:"id" query:"id"`
}

type DeleteVouchersRequest struct {
	IDs []string `json:"ids" query:"ids"`
}
