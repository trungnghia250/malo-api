package delivery

import "github.com/gofiber/fiber/v2"

func (p *ProductHandler) InternalProductAPIRoute(router fiber.Router) {
	router.Get("/product", p.GetProduct)
	router.Delete("/product", p.DeleteProduct)
	router.Post("/product", p.CreateProduct)
	router.Put("/product", p.UpdateProduct)
	router.Get("/product/list", p.ListProduct)
}
