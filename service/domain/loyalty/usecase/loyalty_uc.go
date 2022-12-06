package usecase

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"github.com/trungnghia250/malo-api/service/repo"
)

type ILoyaltyUseCase interface {
	//Gift
	GetGiftByID(ctx *fiber.Ctx, ID string) (*model.Gift, error)
	DeleteGiftsByID(ctx *fiber.Ctx, IDs []string) error
	ListGift(ctx *fiber.Ctx, req dto.ListGiftRequest) ([]model.Gift, error)
	CreateGift(ctx *fiber.Ctx, data *model.Gift) (*model.Gift, error)
	UpdateGift(ctx *fiber.Ctx, data *model.Gift) (*model.Gift, error)

	//Redeem
	GetRedeemByID(ctx *fiber.Ctx, ID string) (*model.RewardRedeem, error)
	DeleteRedeemsByID(ctx *fiber.Ctx, IDs []string) error
	ListRedeem(ctx *fiber.Ctx, req dto.ListRewardRedeemRequest) ([]model.RewardRedeem, error)
	CreateRedeem(ctx *fiber.Ctx, data *model.RewardRedeem) (*model.RewardRedeem, error)
	UpdateRedeem(ctx *fiber.Ctx, data *model.RewardRedeem) (*model.RewardRedeem, error)
}

type loyaltyUseCase struct {
	repo repo.IRepo
}

func NewLoyaltyUseCase(repo repo.IRepo) ILoyaltyUseCase {
	return &loyaltyUseCase{
		repo: repo,
	}
}
