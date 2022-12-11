package delivery

import (
	"github.com/gofiber/fiber/v2"
)

func (c *CampaignHandler) InternalCampaignAPIRoute(router fiber.Router) {
	router.Get("/campaign", c.GetCampaign)
	router.Delete("/campaign", c.DeleteCampaign)
	router.Post("/campaign", c.CreateCampaign)
	router.Put("/campaign", c.UpdateCampaign)
	router.Get("/campaign/list", c.ListCampaign)

	router.Put("/campaign/cancel", c.CancelScheduleCampaign)
}
