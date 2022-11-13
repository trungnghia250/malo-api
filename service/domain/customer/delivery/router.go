package delivery

import (
	"github.com/gofiber/fiber/v2"
)

func (c *CustomerHandler) InternalCustomerAPIRoute(router fiber.Router) {
	router.Get("/customer", c.GetCustomer)
	router.Delete("/customer", c.DeleteCustomer)
	router.Post("/customer", c.CreateCustomer)
	router.Put("/customer", c.UpdateCustomer)
	router.Get("/customer/list", c.ListCustomer)
	router.Put("/customer/add_tags", c.UpdateList)
	router.Get("/customer/export", c.ExportCustomer)
}

func (c *CustomerHandler) ManagementCustomerAPIRoute(router fiber.Router) {
}
