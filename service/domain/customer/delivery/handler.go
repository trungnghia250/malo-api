package delivery

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
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

func (c *CustomerHandler) ListCustomer(ctx *fiber.Ctx) error {
	req := new(dto.ListCustomerRequest)
	if err := ctx.QueryParser(req); err != nil {
		return err
	}

	customers, err := c.customerUseCase.ListCustomer(ctx, *req)
	if err != nil {
		return err
	}

	count := int32(0)
	if len(customers) > 0 {
		count = customers[0].TotalCount
	}

	response := dto.ListCustomerResponse{
		Count: count,
		Data:  customers,
	}

	return ctx.JSON(response)
}

func (c *CustomerHandler) DeleteCustomer(ctx *fiber.Ctx) error {
	req := new(dto.DeleteCustomersRequest)
	if err := ctx.QueryParser(req); err != nil {
		return err
	}

	err := c.customerUseCase.DeleteCustomersByID(ctx, req.CustomerIDs)
	if err != nil {
		return err
	}

	return ctx.JSON(dto.DefaultResponse{
		StatusCode: fiber.StatusOK,
	})
}

func (c *CustomerHandler) UpdateCustomer(ctx *fiber.Ctx) error {
	reqToken := ctx.GetReqHeaders()["X-Access-Token"]
	if reqToken == "" {
		return errors.New("token is required")
	}
	token, err := jwt.Parse(reqToken, nil)
	if token == nil {
		return errors.New("token not valid")
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	req := new(dto.Customer)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	req.ModifiedBy = claims["noc"].(string)
	customer, err := c.customerUseCase.UpdateCustomer(ctx, req)
	if err != nil {
		return err
	}

	return ctx.JSON(customer)
}

func (c *CustomerHandler) CreateCustomer(ctx *fiber.Ctx) error {
	req := new(dto.Customer)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	customer, err := c.customerUseCase.CreateCustomer(ctx, req)
	if err != nil {
		return err
	}

	return ctx.JSON(customer)
}
