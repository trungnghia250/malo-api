package delivery

import (
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

	return ctx.JSON(orders)
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
	req := new(dto.Order)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

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
