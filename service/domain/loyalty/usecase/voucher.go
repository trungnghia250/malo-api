package usecase

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"time"
)

func (l *loyaltyUseCase) GetVoucherByID(ctx *fiber.Ctx, ID string) (*model.Voucher, error) {
	voucher, err := l.repo.NewVoucherRepo().GetVoucherByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	return voucher, nil
}

func (l *loyaltyUseCase) DeleteVouchersByID(ctx *fiber.Ctx, IDs []string) error {
	err := l.repo.NewVoucherRepo().DeleteVouchersByID(ctx, IDs)

	return err
}

func (l *loyaltyUseCase) ListVouchers(ctx *fiber.Ctx, req dto.ListVoucherRequest) ([]model.Voucher, error) {
	vouchers, err := l.repo.NewVoucherRepo().ListVoucher(ctx, req)
	if err != nil {
		return nil, err
	}
	return vouchers, nil
}

func (l *loyaltyUseCase) CreateVoucher(ctx *fiber.Ctx, data *model.Voucher) (*model.Voucher, error) {
	data.CreatedAt = time.Now()

	err := l.repo.NewVoucherRepo().CreateVoucher(ctx, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (l *loyaltyUseCase) UpdateVoucher(ctx *fiber.Ctx, data *model.Voucher) (*model.Voucher, error) {
	data.ModifiedAt = time.Now()
	err := l.repo.NewVoucherRepo().UpdateVoucherByID(ctx, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
