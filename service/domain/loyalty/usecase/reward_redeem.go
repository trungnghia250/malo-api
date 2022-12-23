package usecase

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"time"
)

func (l *loyaltyUseCase) GetRedeemByID(ctx *fiber.Ctx, ID string) (*model.RewardRedeem, error) {
	redeem, err := l.repo.NewRewardRedeem().GetRedeemByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	return redeem, nil
}

func (l *loyaltyUseCase) DeleteRedeemsByID(ctx *fiber.Ctx, IDs []string) error {
	err := l.repo.NewRewardRedeem().DeleteRedeemsByID(ctx, IDs)

	return err
}

func (l *loyaltyUseCase) ListRedeem(ctx *fiber.Ctx, req dto.ListRewardRedeemRequest) ([]model.RewardRedeem, error) {
	redeems, err := l.repo.NewRewardRedeem().ListRedeems(ctx, req)
	if err != nil {
		return nil, err
	}
	return redeems, nil
}

func (l *loyaltyUseCase) CreateRedeem(ctx *fiber.Ctx, data *model.RewardRedeem) (*model.RewardRedeem, error) {
	customer, err := l.repo.NewCustomerRepo().GetCustomerByPhone(ctx, data.Phone)
	if customer.RewardPoint < data.RewardPoint {
		return nil, errors.New("Khách hàng không đủ điểm để đổi phần thưởng này")
	}

	redeemID, err := l.repo.NewCounterRepo().GetSequenceNextValue(ctx, "redeem_id")
	if err != nil {
		return nil, err
	}
	data.ID = fmt.Sprintf("REDEEM%d", redeemID)

	switch data.RewardType {
	case "gift":
		gift, _ := l.repo.NewGiftRepo().GetGiftByID(ctx, "GIFT1")
		if gift.StockAmount == 0 {
			return nil, errors.New("Số lượng phần thưởng đã hết")
		}
		if gift.Status != "ACTIVE" {
			return nil, errors.New("Phần thưởng không được áp dụng")
		}
		if gift.StartAt > int32(time.Now().Unix()) || gift.ExpireAt < int32(time.Now().Unix()) {
			return nil, errors.New("Phần thưởng không áp dụng trong thời gian này")
		}
		err = l.repo.NewCustomerRepo().UpdateCustomerByID(ctx, &dto.Customer{CustomerID: customer.CustomerID, RewardPoint: customer.RewardPoint - data.RewardPoint})

		err = l.repo.NewHistoryPointRepo().CreateHistoryPoint(ctx, &model.HistoryPoint{
			CustomerID:    customer.CustomerID,
			CustomerName:  customer.CustomerName,
			CustomerPhone: customer.PhoneNumber,
			RewardPoint:   data.RewardPoint,
			Type:          "-",
			GiftID:        data.GiftID,
			RedeemID:      fmt.Sprintf("REDEEM%d", redeemID),
			Content:       fmt.Sprintf("Dùng %d điểm đổi quà #%s", data.RewardPoint, data.GiftID),
			CreatedAt:     time.Now(),
		})

		err = l.repo.NewGiftRepo().UpdateGiftByID(ctx, &model.Gift{
			ID:          data.GiftID,
			StockAmount: gift.StockAmount - 1,
			UsedAmount:  gift.StockAmount + 1,
			ModifiedAt:  time.Now(),
			ModifiedBy:  data.ModifiedBy,
		})
	case "discount":
		err = l.repo.NewCustomerRepo().UpdateCustomerByID(ctx, &dto.Customer{CustomerID: customer.CustomerID, RewardPoint: customer.RewardPoint - data.RewardPoint})

		err = l.repo.NewHistoryPointRepo().CreateHistoryPoint(ctx, &model.HistoryPoint{
			CustomerID:    customer.CustomerID,
			CustomerName:  customer.CustomerName,
			CustomerPhone: customer.PhoneNumber,
			RewardPoint:   data.RewardPoint,
			Type:          "-",
			OrderID:       data.OrderID,
			RedeemID:      fmt.Sprintf("REDEEM%d", redeemID),
			Content:       fmt.Sprintf("Dùng %d điểm đổi chiết khấu #%s cho đơn hàng %s", data.RewardPoint, data.GiftID, data.OrderID),
			CreatedAt:     time.Now(),
		})
	}

	data.CreatedAt = time.Now()
	data.ID = fmt.Sprintf("REDEEM%d", redeemID)
	data.RedeemDate = time.Now()
	err = l.repo.NewRewardRedeem().CreateRedeem(ctx, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (l *loyaltyUseCase) UpdateRedeem(ctx *fiber.Ctx, data *model.RewardRedeem) (*model.RewardRedeem, error) {
	data.ModifiedAt = time.Now()
	err := l.repo.NewRewardRedeem().UpdateRedeemByID(ctx, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
