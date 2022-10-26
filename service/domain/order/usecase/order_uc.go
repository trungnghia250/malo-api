package usecase

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"github.com/trungnghia250/malo-api/service/repo"
	"time"
)

type IOrderUseCase interface {
	GetOrderByID(ctx *fiber.Ctx, orderID string) (*model.Order, error)
	DeleteOrderByID(ctx *fiber.Ctx, orderID string) error
	ListOrder(ctx *fiber.Ctx, limit, offset int32) ([]model.Order, error)
	CreateOrder(ctx *fiber.Ctx, data *dto.Order) (*model.Order, error)
	UpdateOrder(ctx *fiber.Ctx, data *dto.Order) (*model.Order, error)
	CountOrder(ctx *fiber.Ctx) (int32, error)
}

type orderUseCase struct {
	repo repo.IRepo
}

func NewOrderUseCase(repo repo.IRepo) IOrderUseCase {
	return &orderUseCase{
		repo: repo,
	}
}

func (o *orderUseCase) GetOrderByID(ctx *fiber.Ctx, orderID string) (*model.Order, error) {
	order, err := o.repo.NewOrderRepo().GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (o *orderUseCase) DeleteOrderByID(ctx *fiber.Ctx, orderID string) error {
	err := o.repo.NewOrderRepo().DeleteOrderByID(ctx, orderID)

	return err
}

func (o *orderUseCase) ListOrder(ctx *fiber.Ctx, limit, offset int32) ([]model.Order, error) {
	orders, err := o.repo.NewOrderRepo().ListOrder(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (o *orderUseCase) CreateOrder(ctx *fiber.Ctx, data *dto.Order) (*model.Order, error) {
	err := o.repo.NewOrderRepo().CreateOrder(ctx, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (o *orderUseCase) UpdateOrder(ctx *fiber.Ctx, data *dto.Order) (*model.Order, error) {
	data.ModifiedAt = time.Now()
	err := o.repo.NewOrderRepo().UpdateOrderByID(ctx, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (o *orderUseCase) CountOrder(ctx *fiber.Ctx) (int32, error) {
	value, err := o.repo.NewOrderRepo().CountOrder(ctx)
	if err != nil {
		return 0, err
	}

	return value, nil
}
