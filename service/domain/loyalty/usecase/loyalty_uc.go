package usecase

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"github.com/trungnghia250/malo-api/service/repo"
)

type ILoyaltyUseCase interface {
	//Gift
	GetProductByID(ctx *fiber.Ctx, ID string) (*model.Gift, error)
	DeleteGiftsByID(ctx *fiber.Ctx, IDs []string) error
	ListGift(ctx *fiber.Ctx, req dto.ListGiftRequest) ([]model.Gift, error)
	CreateGift(ctx *fiber.Ctx, data *model.Gift) (*model.Gift, error)
	UpdateGift(ctx *fiber.Ctx, data *model.Gift) (*model.Gift, error)
}

type loyaltyUseCase struct {
	repo repo.IRepo
}

func NewLoyaltyUseCase(repo repo.IRepo) ILoyaltyUseCase {
	return &loyaltyUseCase{
		repo: repo,
	}
}
