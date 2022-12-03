package delivery

import "github.com/gofiber/fiber/v2"

func (t *TemplateHandler) InternalTemplateAPIRoute(router fiber.Router) {
	router.Get("/template", t.GetTemplate)
	router.Delete("/template", t.DeleteTemplates)
	router.Post("/template", t.CreateTemplate)
	router.Put("/template", t.UpdateTemplate)
	router.Get("/template/list", t.ListTemplate)
}
