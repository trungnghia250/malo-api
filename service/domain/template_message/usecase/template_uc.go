package usecase

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"github.com/trungnghia250/malo-api/service/repo"
	"time"
)

type ITemplateUseCase interface {
	GetTemplateByID(ctx *fiber.Ctx, ID string) (*model.Template, error)
	DeleteTemplateByIDs(ctx *fiber.Ctx, IDs []string) error
	ListTemplate(ctx *fiber.Ctx, req dto.ListTemplateRequest) ([]model.Template, error)
	CreateTemplate(ctx *fiber.Ctx, data *model.Template) (*model.Template, error)
	UpdateTempalte(ctx *fiber.Ctx, data *model.Template) (*model.Template, error)
}

type templateUseCase struct {
	repo repo.IRepo
}

func NewTemplateUseCase(repo repo.IRepo) ITemplateUseCase {
	return &templateUseCase{
		repo: repo,
	}
}

func (t *templateUseCase) GetTemplateByID(ctx *fiber.Ctx, ID string) (*model.Template, error) {
	template, err := t.repo.NewTemplateRepo().GetTemplateByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	return template, nil
}

func (t *templateUseCase) DeleteTemplateByIDs(ctx *fiber.Ctx, IDs []string) error {
	err := t.repo.NewTemplateRepo().DeleteTemplateByID(ctx, IDs)

	return err
}

func (t *templateUseCase) ListTemplate(ctx *fiber.Ctx, req dto.ListTemplateRequest) ([]model.Template, error) {
	templates, err := t.repo.NewTemplateRepo().ListTemplate(ctx, req)
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func (t *templateUseCase) CreateTemplate(ctx *fiber.Ctx, data *model.Template) (*model.Template, error) {
	data.CreatedAt = time.Now()
	templateID, err := t.repo.NewCounterRepo().GetSequenceNextValue(ctx, "template_id")
	if err != nil {
		return nil, err
	}
	data.ID = fmt.Sprintf("MT%d", templateID)
	err = t.repo.NewTemplateRepo().CreateTemplate(ctx, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *templateUseCase) UpdateTempalte(ctx *fiber.Ctx, data *model.Template) (*model.Template, error) {
	data.ModifiedAt = time.Now()
	err := t.repo.NewTemplateRepo().UpdateTemplateByID(ctx, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
