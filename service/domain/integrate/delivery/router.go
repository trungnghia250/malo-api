package delivery

import "github.com/gofiber/fiber/v2"

func (i *IntegrateHandler) InternalIntegrateAPIRoute(router fiber.Router) {
	router.Get("/integrate/partner_config", i.GetPartnerConfig)
	router.Post("/integrate/partner_config", i.CreateConfigPartner)
	router.Put("/integrate/partner_config", i.UpdateConfigPartner)
}
