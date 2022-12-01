package delivery

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/trungnghia250/malo-api/service/domain/customer_group/usecase"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
)

type CustomerGroupHandler struct {
	customerGroupUseCase usecase.ICustomerGroupUseCase
}

func NewCustomerGroupHandler(customerGroupUseCase usecase.ICustomerGroupUseCase) *CustomerGroupHandler {
	return &CustomerGroupHandler{
		customerGroupUseCase: customerGroupUseCase,
	}
}

func (c *CustomerGroupHandler) GetCustomerGroup(ctx *fiber.Ctx) error {
	req := new(dto.GetCustomerGroupByIDRequest)
	if err := ctx.QueryParser(req); err != nil {
		return err
	}

	customerGroup, err := c.customerGroupUseCase.GetCustomerGroupByID(ctx, req.ID)
	if err != nil {
		return err
	}

	return ctx.JSON(customerGroup)
}

func (c *CustomerGroupHandler) ListCustomerGroup(ctx *fiber.Ctx) error {
	req := new(dto.ListCustomerGroupRequest)
	if err := ctx.QueryParser(req); err != nil {
		return err
	}

	groups, err := c.customerGroupUseCase.ListCustomerGroup(ctx, *req)
	if err != nil {
		return err
	}

	count := int32(0)
	if len(groups) > 0 {
		count = groups[0].TotalCount
	}

	response := dto.ListCustomerGroupResponse{
		Count: count,
		Data:  groups,
	}

	return ctx.JSON(response)
}

func (c *CustomerGroupHandler) DeleteCustomerGroup(ctx *fiber.Ctx) error {
	req := new(dto.DeleteCustomerGroupsRequest)
	if err := ctx.QueryParser(req); err != nil {
		return err
	}

	err := c.customerGroupUseCase.DeleteCustomerGroupsByID(ctx, req.IDs)
	if err != nil {
		return err
	}

	return ctx.JSON(dto.DefaultResponse{
		StatusCode: fiber.StatusOK,
	})
}

func (c *CustomerGroupHandler) UpdateCustomerGroup(ctx *fiber.Ctx) error {
	reqToken := ctx.GetReqHeaders()["X-Access-Token"]
	if reqToken == "" {
		return errors.New("token is required")
	}
	token, err := jwt.Parse(reqToken, nil)
	if token == nil {
		return errors.New("token not valid")
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	req := new(model.CustomerGroup)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	req.ModifiedBy = claims["noc"].(string)
	group, err := c.customerGroupUseCase.UpdateCustomerGroup(ctx, req)
	if err != nil {
		return err
	}

	return ctx.JSON(group)
}

func (c *CustomerGroupHandler) CreateCustomerGroup(ctx *fiber.Ctx) error {
	req := new(model.CustomerGroup)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	group, err := c.customerGroupUseCase.CreateCustomerGroup(ctx, req)
	if err != nil {
		return err
	}

	return ctx.JSON(group)
}
