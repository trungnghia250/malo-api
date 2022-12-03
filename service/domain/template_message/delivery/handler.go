package delivery

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/trungnghia250/malo-api/service/domain/template_message/usecase"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
)

type TemplateHandler struct {
	templateUseCase usecase.ITemplateUseCase
}

func NewTemplateHandler(templateUseCase usecase.ITemplateUseCase) *TemplateHandler {
	return &TemplateHandler{
		templateUseCase: templateUseCase,
	}
}

func (t *TemplateHandler) GetTemplate(ctx *fiber.Ctx) error {
	req := new(dto.GetTemplateByIDRequest)
	if err := ctx.QueryParser(req); err != nil {
		return err
	}

	template, err := t.templateUseCase.GetTemplateByID(ctx, req.TemplateID)
	if err != nil {
		return err
	}

	return ctx.JSON(template)
}

func (t *TemplateHandler) ListTemplate(ctx *fiber.Ctx) error {
	req := new(dto.ListTemplateRequest)
	if err := ctx.QueryParser(req); err != nil {
		return err
	}

	templates, err := t.templateUseCase.ListTemplate(ctx, *req)
	if err != nil {
		return err
	}

	count := int32(0)
	if len(templates) > 0 {
		count = templates[0].TotalCount
	}

	response := dto.ListTemplateResponse{
		Count: count,
		Data:  templates,
	}

	return ctx.JSON(response)
}

func (t *TemplateHandler) DeleteTemplates(ctx *fiber.Ctx) error {
	req := new(dto.DeleteTemplatesRequest)
	if err := ctx.QueryParser(req); err != nil {
		return err
	}

	err := t.templateUseCase.DeleteTemplateByIDs(ctx, req.IDs)
	if err != nil {
		return err
	}

	return ctx.JSON(dto.DefaultResponse{
		StatusCode: fiber.StatusOK,
	})
}

func (t *TemplateHandler) UpdateTemplate(ctx *fiber.Ctx) error {
	reqToken := ctx.GetReqHeaders()["X-Access-Token"]
	if reqToken == "" {
		return errors.New("token is required")
	}
	token, err := jwt.Parse(reqToken, nil)
	if token == nil {
		return errors.New("token not valid")
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	req := new(model.Template)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}
	req.ModifiedBy = claims["noc"].(string)
	template, err := t.templateUseCase.UpdateTempalte(ctx, req)
	if err != nil {
		return err
	}

	return ctx.JSON(template)
}

func (t *TemplateHandler) CreateTemplate(ctx *fiber.Ctx) error {
	req := new(model.Template)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}
	reqToken := ctx.GetReqHeaders()["X-Access-Token"]
	if reqToken == "" {
		return errors.New("token is required")
	}
	token, err := jwt.Parse(reqToken, nil)
	if token == nil {
		return errors.New("token not valid")
	}
	claims, _ := token.Claims.(jwt.MapClaims)

	req.CreatedBy = claims["noc"].(string)
	template, err := t.templateUseCase.CreateTemplate(ctx, req)
	if err != nil {
		return err
	}

	return ctx.JSON(template)
}
