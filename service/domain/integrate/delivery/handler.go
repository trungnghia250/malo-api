package delivery

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/domain/integrate/usecase"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
)

type IntegrateHandler struct {
	integrateUseCase usecase.IIntegrateUseCase
}

func NewIntegrateHandler(integrateUseCase usecase.IIntegrateUseCase) *IntegrateHandler {
	return &IntegrateHandler{
		integrateUseCase: integrateUseCase,
	}
}

func (i *IntegrateHandler) GetPartnerConfig(ctx *fiber.Ctx) error {
	req := new(dto.GetPartnerConfigByIDRequest)
	if err := ctx.QueryParser(req); err != nil {
		return err
	}

	config, err := i.integrateUseCase.GetPartnerConfig(ctx, req.PartnerCode)
	if err != nil {
		return err
	}

	return ctx.JSON(config)
}

func (i *IntegrateHandler) UpdateConfigPartner(ctx *fiber.Ctx) error {
	reqToken := ctx.GetReqHeaders()["X-Access-Token"]
	if reqToken == "" {
		return errors.New("token is required")
	}
	token, err := jwt.Parse(reqToken, nil)
	if token == nil {
		return errors.New("token not valid")
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	req := new(model.PartnerConfig)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	req.ModifiedBy = claims["noc"].(string)
	config, err := i.integrateUseCase.UpdatePartnerConfig(ctx, req)
	if err != nil {
		return err
	}

	return ctx.JSON(config)
}

func (i *IntegrateHandler) CreateConfigPartner(ctx *fiber.Ctx) error {
	req := new(model.PartnerConfig)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	config, err := i.integrateUseCase.CreatePartnerConfig(ctx, req)
	if err != nil {
		return err
	}

	return ctx.JSON(config)
}
