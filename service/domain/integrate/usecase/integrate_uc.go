package usecase

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/repo"
	"strings"
	"time"
)

type IIntegrateUseCase interface {
	GetPartnerConfig(ctx *fiber.Ctx, partner string) (*model.PartnerConfig, error)
	CreatePartnerConfig(ctx *fiber.Ctx, data *model.PartnerConfig) (*model.PartnerConfig, error)
	UpdatePartnerConfig(ctx *fiber.Ctx, data *model.PartnerConfig) (*model.PartnerConfig, error)
}

type integrateUseCase struct {
	repo repo.IRepo
}

func NewIntegrateUseCase(repo repo.IRepo) IIntegrateUseCase {
	return &integrateUseCase{
		repo: repo,
	}
}

func (i *integrateUseCase) GetPartnerConfig(ctx *fiber.Ctx, partner string) (*model.PartnerConfig, error) {
	partner = strings.ToUpper(partner)
	partnerConfig, err := i.repo.NewPartnerRepo().GetPartnerConfig(ctx, partner)
	if err != nil {
		return nil, err
	}
	return partnerConfig, nil
}

func (i *integrateUseCase) CreatePartnerConfig(ctx *fiber.Ctx, data *model.PartnerConfig) (*model.PartnerConfig, error) {
	data.CreatedAt = time.Now()
	data.ModifiedAt = time.Now()

	err := i.repo.NewPartnerRepo().CreatePartnerConfig(ctx, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (i *integrateUseCase) UpdatePartnerConfig(ctx *fiber.Ctx, data *model.PartnerConfig) (*model.PartnerConfig, error) {
	data.ModifiedAt = time.Now()

	err := i.repo.NewPartnerRepo().UpdatePartnerByID(ctx, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
