package delivery

import "github.com/gofiber/fiber/v2"

func (l *LoyaltyHandler) InternalLoyaltyAPIRoute(router fiber.Router) {
	//gift
	router.Get("/loyalty/gift", l.GetGift)
	router.Delete("/loyalty/gift", l.DeleteGift)
	router.Post("/loyalty/gift", l.CreateGift)
	router.Put("/loyalty/gift", l.UpdateGift)
	router.Get("/loyalty/gift/list", l.ListGift)

	//reward_redeem
	router.Get("/loyalty/redeem", l.GetRedeem)
	router.Delete("/loyalty/redeem", l.DeleteRedeem)
	router.Post("/loyalty/redeem", l.CreateRedeem)
	router.Put("/loyalty/redeem", l.UpdateRedeem)
	router.Get("/loyalty/redeem/list", l.ListRedeem)

	//voucher
	router.Get("/loyalty/voucher", l.GetVoucher)
	router.Delete("/loyalty/voucher", l.DeleteVoucher)
	router.Post("/loyalty/voucher", l.CreateVoucher)
	router.Put("/loyalty/voucher", l.UpdateVoucher)
	router.Get("/loyalty/voucher/list", l.ListVoucher)
	router.Get("loyalty/voucher/validate", l.ValidateVoucher)

	//voucher_usage
	router.Get("/loyalty/voucher_usage", l.GetVoucherUsage)
	router.Delete("/loyalty/voucher_usage", l.DeleteVoucherUsage)
	router.Post("/loyalty/voucher_usage", l.CreateVoucherUsage)
	router.Put("/loyalty/voucher_usage", l.UpdateVoucherUsage)
	router.Get("/loyalty/voucher_usage/list", l.ListVoucherUsage)

	//earn_point and rank
	router.Get("/loyalty/config", l.GetLoyaltyConfig)
	router.Put("/loyalty/config", l.UpdateLoyaltyConfig)
}
