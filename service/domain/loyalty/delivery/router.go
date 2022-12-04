package delivery

import "github.com/gofiber/fiber/v2"

func (l *LoyaltyHandler) InternalLoyaltyAPIRoute(router fiber.Router) {
	router.Get("/loyalty/gift", l.GetGift)
	router.Delete("/loyalty/gift", l.DeleteGift)
	router.Post("/loyalty/gift", l.CreateGift)
	router.Put("/loyalty/gift", l.UpdateGift)
	router.Get("/loyalty/gift/list", l.ListGift)
}
