package usecase

import (
	"errors"
	"fmt"
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

func (l *loyaltyUseCase) ListVouchers(ctx *fiber.Ctx, req dto.ListVoucherRequest) (resp dto.ListVoucherResponse, err error) {
	vouchers, err := l.repo.NewVoucherRepo().ListVoucher(ctx, req)
	if err != nil {
		return dto.ListVoucherResponse{}, err
	}
	count := int32(0)
	if len(vouchers) > 0 {
		count = vouchers[0].TotalCount
	}
	resp.Count = count

	for _, voucher := range vouchers {
		groups, _ := l.repo.NewCustomerGroupRepo().ListCustomerGroup(ctx, dto.ListCustomerGroupRequest{
			Limit: int32(len(voucher.GroupIDs)),
			IDs:   voucher.GroupIDs,
		})

		var groupNames []string
		for _, group := range groups {
			groupNames = append(groupNames, group.GroupName)
		}

		resp.Data = append(resp.Data, dto.VoucherInGroup{
			Code:           voucher.ID,
			GroupNames:     groupNames,
			DiscountAmount: voucher.DiscountAmount,
			StartAt:        voucher.StartAt,
			ExpireAt:       voucher.ExpireAt,
			Note:           voucher.Note,
			CreatedAt:      voucher.CreatedAt,
			Status:         voucher.Status,
		})
	}
	return resp, nil
}

func (l *loyaltyUseCase) CreateVoucher(ctx *fiber.Ctx, data *model.Voucher) (*model.Voucher, error) {
	data.CreatedAt = time.Now()

	_, err := l.GetVoucherByID(ctx, data.ID)
	if err.Error() != "mongo: no documents in result" {
		return nil, errors.New("Voucher Code đã tồn tại trong hệ thống")
	}

	err = l.repo.NewVoucherRepo().CreateVoucher(ctx, data)
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

func (l *loyaltyUseCase) ValidateVoucher(ctx *fiber.Ctx, req dto.ValidateVoucherRequest) (dto.ValidateVoucherResponse, error) {
	var res dto.ValidateVoucherResponse
	customer, err := l.repo.NewCustomerRepo().GetCustomerByPhone(ctx, req.Phone)
	if err != nil {
		return dto.ValidateVoucherResponse{}, err
	}

	res.CustomerDetail = dto.CustomerDetail{
		Name:         customer.CustomerName,
		Phone:        customer.PhoneNumber,
		Address:      customer.Address,
		RewardPoint:  customer.RewardPoint,
		Gender:       customer.Gender,
		Email:        customer.Email,
		CustomerType: customer.CustomerType,
	}

	groups, err := l.repo.NewCustomerGroupRepo().ListCustomerGroupByCustomerID(ctx, customer.CustomerID)
	var groupIDs []string
	for _, group := range groups {
		groupIDs = append(groupIDs, group.ID)
	}
	vouchers, err := l.repo.NewVoucherRepo().ListValidateVoucherByGroupIDs(ctx, groupIDs)
	if len(req.Code) > 0 {
		res.CheckVoucherMessage = fmt.Sprintf("Mã %s không còn hoạt động", req.Code)
	}
	for _, voucher := range vouchers {
		countUsed, _ := l.repo.NewVoucherUsageRepo().CountCustomerUseVoucher(ctx, customer.PhoneNumber, voucher.ID)
		if voucher.ID == req.Code {
			res.CheckVoucherMessage = fmt.Sprintf("Mã %s khách hàng hết số lần sử dụng", req.Code)
		}
		if countUsed < voucher.LimitPerCustomer {
			res.Vouchers = append(res.Vouchers, dto.VoucherDetail{
				Code:             voucher.ID,
				DiscountAmount:   voucher.DiscountAmount,
				MinOrderAmount:   voucher.MinOrderAmount,
				StartAt:          voucher.StartAt,
				ExpireAt:         voucher.ExpireAt,
				CustomerUsed:     countUsed,
				LimitPerCustomer: voucher.LimitPerCustomer,
			})
			if voucher.ID == req.Code {
				res.CheckVoucherMessage = fmt.Sprintf("Mã %s khách hàng có thể sử dụng", req.Code)
			}
		}

	}

	gifts, _ := l.repo.NewGiftRepo().ListGiftValidateCustomer(ctx, customer.RewardPoint)
	for _, gift := range gifts {
		res.Gifts = append(res.Gifts, dto.GiftDetail{
			ID:          gift.ID,
			Name:        gift.Name,
			URL:         gift.ImageURL,
			Price:       gift.Price,
			RewardPoint: gift.RewardPoint,
			StockAmount: gift.StockAmount,
		})
	}

	return res, nil
}
