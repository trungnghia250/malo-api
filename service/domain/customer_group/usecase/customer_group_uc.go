package usecase

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"github.com/trungnghia250/malo-api/service/repo"
	"time"
)

type ICustomerGroupUseCase interface {
	GetCustomerGroupByID(ctx *fiber.Ctx, customerGroupID string) (dto.GetCustomerGroupResponse, error)
	DeleteCustomerGroupsByID(ctx *fiber.Ctx, customerGroupIDs []string) error
	ListCustomerGroup(ctx *fiber.Ctx, req dto.ListCustomerGroupRequest) ([]model.CustomerGroup, error)
	CreateCustomerGroup(ctx *fiber.Ctx, data *model.CustomerGroup) (*model.CustomerGroup, error)
	UpdateCustomerGroup(ctx *fiber.Ctx, data *model.CustomerGroup) (*model.CustomerGroup, error)
}

type customerGroupUseCase struct {
	repo repo.IRepo
}

func NewCustomerGroupUseCase(repo repo.IRepo) ICustomerGroupUseCase {
	return &customerGroupUseCase{
		repo: repo,
	}
}

func (c *customerGroupUseCase) GetCustomerGroupByID(ctx *fiber.Ctx, customerGroupID string) (dto.GetCustomerGroupResponse, error) {
	customerGroup, err := c.repo.NewCustomerGroupRepo().GetCustomerGroupByID(ctx, customerGroupID)
	if err != nil {
		return dto.GetCustomerGroupResponse{}, err
	}
	customers, err := c.repo.NewCustomerRepo().ListCustomer(ctx, dto.ListCustomerRequest{
		Limit:       int32(len(customerGroup.CustomerIDs)),
		CustomerIDs: customerGroup.CustomerIDs,
	})

	return dto.GetCustomerGroupResponse{
		ID:         customerGroup.ID,
		GroupName:  customerGroup.GroupName,
		Note:       customerGroup.Note,
		CreatedAt:  customerGroup.CreatedAt,
		ModifiedAt: customerGroup.ModifiedAt,
		ModifiedBy: customerGroup.ModifiedBy,
		Customers:  customers,
	}, nil
}

func (c *customerGroupUseCase) DeleteCustomerGroupsByID(ctx *fiber.Ctx, customerGroupIDs []string) error {
	err := c.repo.NewCustomerGroupRepo().DeleteCustomerGroupByID(ctx, customerGroupIDs)

	return err
}

func (c *customerGroupUseCase) ListCustomerGroup(ctx *fiber.Ctx, req dto.ListCustomerGroupRequest) ([]model.CustomerGroup, error) {
	customerGroups, err := c.repo.NewCustomerGroupRepo().ListCustomerGroup(ctx, req)
	if err != nil {
		return nil, err
	}
	return customerGroups, nil
}

func (c *customerGroupUseCase) CreateCustomerGroup(ctx *fiber.Ctx, data *model.CustomerGroup) (*model.CustomerGroup, error) {
	data.CreatedAt = time.Now()
	customerGroupID, err := c.repo.NewCounterRepo().GetSequenceNextValue(ctx, "customer_group_id")
	if err != nil {
		return nil, err
	}
	data.ID = fmt.Sprintf("CG%d", customerGroupID)
	err = c.repo.NewCustomerGroupRepo().CreateCustomerGroup(ctx, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *customerGroupUseCase) UpdateCustomerGroup(ctx *fiber.Ctx, data *model.CustomerGroup) (*model.CustomerGroup, error) {
	data.ModifiedAt = time.Now()

	err := c.repo.NewCustomerGroupRepo().UpdateCustomerGroupByID(ctx, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
