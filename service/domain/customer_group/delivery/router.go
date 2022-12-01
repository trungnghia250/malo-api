package delivery

import (
	"github.com/gofiber/fiber/v2"
)

func (c *CustomerGroupHandler) InternalCustomerGroupAPIRoute(router fiber.Router) {
	router.Get("/customer_group", c.GetCustomerGroup)
	router.Delete("/customer_group", c.DeleteCustomerGroup)
	router.Post("/customer_group", c.CreateCustomerGroup)
	router.Put("/customer_group", c.UpdateCustomerGroup)
	router.Get("/customer_group/list", c.ListCustomerGroup)
}
