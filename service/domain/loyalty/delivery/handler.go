package delivery

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/domain/loyalty/usecase"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
)

type LoyaltyHandler struct {
	loyaltyUseCase usecase.ILoyaltyUseCase
}

func NewLoyaltyHandler(loyaltyUseCase usecase.ILoyaltyUseCase) *LoyaltyHandler {
	return &LoyaltyHandler{
		loyaltyUseCase: loyaltyUseCase,
	}
}

func (l *LoyaltyHandler) GetGift(ctx *fiber.Ctx) error {
	req := new(dto.GetGiftByIDRequest)
	if err := ctx.QueryParser(req); err != nil {
		return err
	}

	gift, err := l.loyaltyUseCase.GetProductByID(ctx, req.GiftID)
	if err != nil {
		return err
	}

	return ctx.JSON(gift)
}

func (l *LoyaltyHandler) ListGift(ctx *fiber.Ctx) error {
	req := new(dto.ListGiftRequest)
	if err := ctx.QueryParser(req); err != nil {
		return err
	}

	gifts, err := l.loyaltyUseCase.ListGift(ctx, *req)
	if err != nil {
		return err
	}

	count := int32(0)
	if len(gifts) > 0 {
		count = gifts[0].TotalCount
	}

	response := dto.ListGiftResponse{
		Count: count,
		Data:  gifts,
	}

	return ctx.JSON(response)
}

func (l *LoyaltyHandler) DeleteGift(ctx *fiber.Ctx) error {
	req := new(dto.DeleteGiftsRequest)
	if err := ctx.QueryParser(req); err != nil {
		return err
	}

	err := l.loyaltyUseCase.DeleteGiftsByID(ctx, req.IDs)
	if err != nil {
		return err
	}

	return ctx.JSON(dto.DefaultResponse{
		StatusCode: fiber.StatusOK,
	})
}

func (l *LoyaltyHandler) UpdateGift(ctx *fiber.Ctx) error {
	reqToken := ctx.GetReqHeaders()["X-Access-Token"]
	if reqToken == "" {
		return errors.New("token is required")
	}
	token, err := jwt.Parse(reqToken, nil)
	if token == nil {
		return errors.New("token not valid")
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	req := new(model.Gift)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}
	req.ModifiedBy = claims["noc"].(string)
	gift, err := l.loyaltyUseCase.UpdateGift(ctx, req)
	if err != nil {
		return err
	}

	return ctx.JSON(gift)
}

func (l *LoyaltyHandler) CreateGift(ctx *fiber.Ctx) error {
	reqToken := ctx.GetReqHeaders()["X-Access-Token"]
	if reqToken == "" {
		return errors.New("token is required")
	}
	token, err := jwt.Parse(reqToken, nil)
	if token == nil {
		return errors.New("token not valid")
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	req := new(model.Gift)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}
	req.CreatedBy = claims["noc"].(string)
	gift, err := l.loyaltyUseCase.CreateGift(ctx, req)
	if err != nil {
		return err
	}

	return ctx.JSON(gift)
}
