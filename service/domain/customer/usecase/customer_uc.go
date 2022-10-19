package usecase

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"github.com/trungnghia250/malo-api/service/repo"
)

type ICustomerUseCase interface {
	GetCustomerByID(ctx *fiber.Ctx, customerID string) (*model.Customer, error)
	DeleteCustomerByID(ctx *fiber.Ctx, customerID string) error
	ListCustomer(ctx *fiber.Ctx, limit, offset int32) ([]model.Customer, error)
	CreateCustomer(ctx *fiber.Ctx, data *dto.Customer) (*model.Customer, error)
	UpdateCustomer(ctx *fiber.Ctx, data *dto.Customer) (*model.Customer, error)
}

type customerUseCase struct {
	repo repo.IRepo
}

func NewCustomerUseCase(repo repo.IRepo) ICustomerUseCase {
	return &customerUseCase{
		repo: repo,
	}
}

func (c *customerUseCase) GetCustomerByID(ctx *fiber.Ctx, customerID string) (*model.Customer, error) {
	customer, err := c.repo.NewCustomerRepo().GetCustomerByID(ctx, customerID)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (c *customerUseCase) DeleteCustomerByID(ctx *fiber.Ctx, customerID string) error {
	err := c.repo.NewCustomerRepo().DeleteCustomerByID(ctx, customerID)

	return err
}

func (c *customerUseCase) ListCustomer(ctx *fiber.Ctx, limit, offset int32) ([]model.Customer, error) {
	customers, err := c.repo.NewCustomerRepo().ListCustomer(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return customers, nil
}

func (c *customerUseCase) CreateCustomer(ctx *fiber.Ctx, data *dto.Customer) (*model.Customer, error) {
	err := c.repo.NewCustomerRepo().CreateCustomer(ctx, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *customerUseCase) UpdateCustomer(ctx *fiber.Ctx, data *dto.Customer) (*model.Customer, error) {
	err := c.repo.NewCustomerRepo().UpdateCustomerByID(ctx, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}