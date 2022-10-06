package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/domain/customer/usecase"
	"github.com/trungnghia250/malo-api/service/model/dto"
)

type CustomerHandler struct {
	customerUseCase usecase.ICustomerUseCase
}

func NewCustomerHandler(customerUseCase usecase.ICustomerUseCase) *CustomerHandler {
	return &CustomerHandler{
		customerUseCase: customerUseCase,
	}
}

func (c *CustomerHandler) GetCustomer(ctx *fiber.Ctx) error {
	req := new(dto.GetCustomerByIDRequest)
	if err := ctx.QueryParser(req); err != nil {
		return err
	}

	customer, err := c.customerUseCase.GetCustomerByID(ctx, req.CustomerID)
	if err != nil {
		return err
	}

	return ctx.JSON(customer)
}
