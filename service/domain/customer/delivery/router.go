package delivery

import (
	"github.com/gofiber/fiber/v2"
)

func (c *CustomerHandler) InternalCustomerAPIRoute(router fiber.Router) {
	router.Get("/customer", c.GetCustomer)
}

func (c *CustomerHandler) ManagementCustomerAPIRoute(router fiber.Router) {
}
