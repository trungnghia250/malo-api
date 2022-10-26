package delivery

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/domain/order/usecase"
	"github.com/trungnghia250/malo-api/service/model/dto"
)

type OrderHandler struct {
	orderUseCase usecase.IOrderUseCase
}

func NewOrderHandler(orderUseCase usecase.IOrderUseCase) *OrderHandler {
	return &OrderHandler{
		orderUseCase: orderUseCase,
	}
}

func (o *OrderHandler) GetOrder(ctx *fiber.Ctx) error {
	req := new(dto.GetOrderByIDRequest)
	if err := ctx.QueryParser(req); err != nil {
		return err
	}

	order, err := o.orderUseCase.GetOrderByID(ctx, req.OrderID)
	if err != nil {
		return err
	}

	return ctx.JSON(order)
}

func (o *OrderHandler) ListOrder(ctx *fiber.Ctx) error {
	req := new(dto.ListCustomerRequest)
	if err := ctx.QueryParser(req); err != nil {
		return err
	}

	orders, err := o.orderUseCase.ListOrder(ctx, req.Limit, req.Offset)
	if err != nil {
		return err
	}

	count, err := o.orderUseCase.CountOrder(ctx)
	if err != nil {
		return err
	}

	response := dto.ListOrderResponse{
		Count: count,
		Data:  orders,
	}

	return ctx.JSON(response)
}

func (o *OrderHandler) DeleteOrder(ctx *fiber.Ctx) error {
	req := new(dto.GetOrderByIDRequest)
	if err := ctx.QueryParser(req); err != nil {
		return err
	}

	err := o.orderUseCase.DeleteOrderByID(ctx, req.OrderID)
	if err != nil {
		return err
	}

	return ctx.JSON(dto.DefaultResponse{
		StatusCode: fiber.StatusOK,
	})
}

func (o *OrderHandler) UpdateOrder(ctx *fiber.Ctx) error {
	reqToken := ctx.GetReqHeaders()["X-Access-Token"]
	if reqToken == "" {
		return errors.New("token is required")
	}
	token, err := jwt.Parse(reqToken, nil)
	if token == nil {
		return errors.New("token not valid")
	}
	claims, _ := token.Claims.(jwt.MapClaims)

	req := new(dto.Order)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}
	req.ModifiedBy = claims["noc"].(string)
	order, err := o.orderUseCase.UpdateOrder(ctx, req)
	if err != nil {
		return err
	}

	return ctx.JSON(order)
}

func (o *OrderHandler) CreateOrder(ctx *fiber.Ctx) error {
	req := new(dto.Order)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	order, err := o.orderUseCase.CreateOrder(ctx, req)
	if err != nil {
		return err
	}

	return ctx.JSON(order)
}
