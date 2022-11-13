package delivery

import (
	"github.com/gofiber/fiber/v2"
)

func (o *OrderHandler) InternalOrderAPIRoute(router fiber.Router) {
	router.Get("/order", o.GetOrder)
	router.Delete("/order", o.DeleteOrder)
	router.Post("/order", o.CreateOrder)
	router.Put("/order", o.UpdateOrder)
	router.Get("/order/list", o.ListOrder)
	router.Post("/order/import", o.ImportOrder)
	router.Get("/order/export", o.ExportOrder)
}
