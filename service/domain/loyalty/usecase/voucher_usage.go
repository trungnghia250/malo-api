package usecase

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"time"
)

func (l *loyaltyUseCase) GetVoucherUsageByID(ctx *fiber.Ctx, ID string) (*model.VoucherUsage, error) {
	usage, err := l.repo.NewVoucherUsageRepo().GetVoucherUsageByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	return usage, nil
}

func (l *loyaltyUseCase) DeleteVoucherUsagesByID(ctx *fiber.Ctx, IDs []string) error {
	err := l.repo.NewVoucherUsageRepo().DeleteVoucherUsagesByID(ctx, IDs)

	return err
}

func (l *loyaltyUseCase) ListVoucherUsages(ctx *fiber.Ctx, req dto.ListVoucherUsageRequest) ([]model.VoucherUsage, error) {
	usages, err := l.repo.NewVoucherUsageRepo().ListVoucherUsages(ctx, req)
	if err != nil {
		return nil, err
	}
	return usages, nil
}

func (l *loyaltyUseCase) CreateVoucherUsage(ctx *fiber.Ctx, data *model.VoucherUsage) (*model.VoucherUsage, error) {
	usageID, err := l.repo.NewCounterRepo().GetSequenceNextValue(ctx, "voucher_usage_id")
	if err != nil {
		return nil, err
	}
	data.ID = fmt.Sprintf("VU%d", usageID)
	data.CreatedAt = time.Now()
	err = l.repo.NewVoucherUsageRepo().CreateVoucherUsage(ctx, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (l *loyaltyUseCase) UpdateVoucherUsage(ctx *fiber.Ctx, data *model.VoucherUsage) (*model.VoucherUsage, error) {
	data.ModifiedAt = time.Now()
	err := l.repo.NewVoucherUsageRepo().UpdateVoucherUsageByID(ctx, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
