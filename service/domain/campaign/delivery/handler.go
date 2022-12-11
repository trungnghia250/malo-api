package delivery

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/trungnghia250/malo-api/service/domain/campaign/usecase"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
)

type CampaignHandler struct {
	campaignUseCase usecase.ICampaignUseCase
}

func NewCampaignHandler(campaignUseCase usecase.ICampaignUseCase) *CampaignHandler {
	return &CampaignHandler{
		campaignUseCase: campaignUseCase,
	}
}

func (c *CampaignHandler) GetCampaign(ctx *fiber.Ctx) error {
	req := new(dto.GetCampaignByIDRequest)
	if err := ctx.QueryParser(req); err != nil {
		return err
	}

	campaign, err := c.campaignUseCase.GetCampaignByID(ctx, req.CampaignID)
	if err != nil {
		return err
	}

	return ctx.JSON(campaign)
}

func (c *CampaignHandler) ListCampaign(ctx *fiber.Ctx) error {
	req := new(dto.ListCampaignRequest)
	if err := ctx.QueryParser(req); err != nil {
		return err
	}

	campaigns, err := c.campaignUseCase.ListCampaign(ctx, *req)
	if err != nil {
		return err
	}

	count := int32(0)
	if len(campaigns) > 0 {
		count = campaigns[0].TotalCount
	}

	response := dto.ListCampaignResponse{
		Count: count,
		Data:  campaigns,
	}

	return ctx.JSON(response)
}

func (c *CampaignHandler) DeleteCampaign(ctx *fiber.Ctx) error {
	req := new(dto.DeleteCampaignsRequest)
	if err := ctx.QueryParser(req); err != nil {
		return err
	}

	err := c.campaignUseCase.DeleteCampaignsByID(ctx, req.CampaignIDs)
	if err != nil {
		return err
	}

	return ctx.JSON(dto.DefaultResponse{
		StatusCode: fiber.StatusOK,
	})
}

func (c *CampaignHandler) UpdateCampaign(ctx *fiber.Ctx) error {
	reqToken := ctx.GetReqHeaders()["X-Access-Token"]
	if reqToken == "" {
		return errors.New("token is required")
	}
	token, err := jwt.Parse(reqToken, nil)
	if token == nil {
		return errors.New("token not valid")
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	req := new(model.Campaign)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	req.ModifiedBy = claims["noc"].(string)
	customer, err := c.campaignUseCase.UpdateCampaign(ctx, req)
	if err != nil {
		return err
	}

	return ctx.JSON(customer)
}

func (c *CampaignHandler) CreateCampaign(ctx *fiber.Ctx) error {
	reqToken := ctx.GetReqHeaders()["X-Access-Token"]
	if reqToken == "" {
		return errors.New("token is required")
	}
	token, err := jwt.Parse(reqToken, nil)
	if token == nil {
		return errors.New("token not valid")
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	req := new(model.Campaign)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	req.CreatedBy = claims["noc"].(string)
	campaign, err := c.campaignUseCase.CreateCampaign(ctx, req)
	if err != nil {
		return err
	}

	return ctx.JSON(campaign)
}

func (c *CampaignHandler) CancelScheduleCampaign(ctx *fiber.Ctx) error {
	reqToken := ctx.GetReqHeaders()["X-Access-Token"]
	if reqToken == "" {
		return errors.New("token is required")
	}
	token, err := jwt.Parse(reqToken, nil)
	if token == nil {
		return errors.New("token not valid")
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	req := new(model.Campaign)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	req.ModifiedBy = claims["noc"].(string)
	customer, err := c.campaignUseCase.CancelCampaign(ctx, req)
	if err != nil {
		return err
	}

	return ctx.JSON(customer)
}
