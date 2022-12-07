package usecase

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/model"
)

func (l *loyaltyUseCase) GetLoyaltyConfig(ctx *fiber.Ctx) (*model.LoyaltyConfig, error) {
	config, err := l.repo.NewPartnerRepo().GetPartnerLoyaltyConfig(ctx)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (l *loyaltyUseCase) UpdateLoyaltyConfig(ctx *fiber.Ctx, data *model.LoyaltyConfig) (*model.LoyaltyConfig, error) {
	err := l.repo.NewPartnerRepo().UpdatePartnerLoyaltyConfig(ctx, data)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
