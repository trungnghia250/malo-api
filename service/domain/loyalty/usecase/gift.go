package usecase

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"time"
)

func (l *loyaltyUseCase) GetGiftByID(ctx *fiber.Ctx, ID string) (*model.Gift, error) {
	gift, err := l.repo.NewGiftRepo().GetGiftByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	return gift, nil
}

func (l *loyaltyUseCase) DeleteGiftsByID(ctx *fiber.Ctx, IDs []string) error {
	err := l.repo.NewGiftRepo().DeleteGiftsByID(ctx, IDs)

	return err
}

func (l *loyaltyUseCase) ListGift(ctx *fiber.Ctx, req dto.ListGiftRequest) ([]model.Gift, error) {
	gifts, err := l.repo.NewGiftRepo().ListGift(ctx, req)
	if err != nil {
		return nil, err
	}
	return gifts, nil
}

func (l *loyaltyUseCase) CreateGift(ctx *fiber.Ctx, data *model.Gift) (*model.Gift, error) {
	data.CreatedAt = time.Now()
	giftID, err := l.repo.NewCounterRepo().GetSequenceNextValue(ctx, "gift_id")
	if err != nil {
		return nil, err
	}
	data.ID = fmt.Sprintf("GIFT%d", giftID)
	err = l.repo.NewGiftRepo().CreateGift(ctx, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (l *loyaltyUseCase) UpdateGift(ctx *fiber.Ctx, data *model.Gift) (*model.Gift, error) {
	data.ModifiedAt = time.Now()
	err := l.repo.NewGiftRepo().UpdateGiftByID(ctx, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
