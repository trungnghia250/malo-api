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

type ListVoucherUsageRequest struct {
	Limit          int32    `json:"limit,omitempty" query:"limit,omitempty"`
	Offset         int32    `json:"offset,omitempty" query:"offset,omitempty"`
	Code           []string `json:"code,omitempty" query:"code,omitempty"`
	CustomerName   []string `json:"customer_name,omitempty" query:"customer_name,omitempty"`
	Phone          []string `json:"phone,omitempty" query:"phone,omitempty"`
	OrderID        []string `json:"order_id,omitempty" query:"order_id,omitempty"`
	DiscountAmount []int32  `json:"discount_amount,omitempty" query:"discount_amount,omitempty"`
	UsageDate      []int32  `json:"usage_date,omitempty" query:"usage_date,omitempty"`
	CreatedAt      []int32  `json:"created_at,omitempty" query:"created_at,omitempty"`
}

type ListVoucherUsageResponse struct {
	Count int32                `json:"count"`
	Data  []model.VoucherUsage `json:"data"`
}

type ValidateVoucherRequest struct {
	Phone string `json:"phone,omitempty" query:"phone,omitempty"`
	Code  string `json:"code,omitempty" query:"code,omitempty"`
}

type ValidateVoucherResponse struct {
	CustomerDetail      CustomerDetail  `json:"customer_detail"`
	Vouchers            []VoucherDetail `json:"vouchers,omitempty"`
	Gifts               []GiftDetail    `json:"gifts,omitempty"`
	CheckVoucherMessage string          `json:"check_voucher_message,omitempty"`
}

type CustomerDetail struct {
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	RewardPoint  int32  `json:"reward_point"`
	Gender       string `json:"gender"`
	Email        string `json:"email"`
	CustomerType string `json:"customer_type"`
}

type VoucherDetail struct {
	Code             string `json:"code"`
	DiscountAmount   int32  `json:"discount_amount"`
	MinOrderAmount   int32  `json:"min_order_amount"`
	StartAt          int32  `json:"start_date"`
	ExpireAt         int32  `json:"expire_at"`
	CustomerUsed     int32  `json:"customer_used"`
	LimitPerCustomer int32  `json:"limit_per_customer,omitempty"`
}

type GiftDetail struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Price       int32  `json:"price"`
	RewardPoint int32  `json:"reward_point"`
	StockAmount int32  `json:"stock_amount"`
}
