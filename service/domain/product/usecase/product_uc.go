package usecase

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"github.com/trungnghia250/malo-api/service/repo"
	"time"
)

type IProductUseCase interface {
	GetProductByID(ctx *fiber.Ctx, productID string) (*model.Product, error)
	DeleteProductByID(ctx *fiber.Ctx, productID string) error
	ListProduct(ctx *fiber.Ctx, req dto.ListProductRequest) ([]model.Product, error)
	CreateProduct(ctx *fiber.Ctx, data *model.Product) (*model.Product, error)
	UpdateProduct(ctx *fiber.Ctx, data *model.Product) (*model.Product, error)
	CountProduct(ctx *fiber.Ctx) (int32, error)
}

type productUseCase struct {
	repo repo.IRepo
}

func NewProductUseCase(repo repo.IRepo) IProductUseCase {
	return &productUseCase{
		repo: repo,
	}
}

func (p *productUseCase) GetProductByID(ctx *fiber.Ctx, productID string) (*model.Product, error) {
	product, err := p.repo.NewProductRepo().GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *productUseCase) DeleteProductByID(ctx *fiber.Ctx, productID string) error {
	err := p.repo.NewProductRepo().DeleteProductByID(ctx, productID)

	return err
}

func (p *productUseCase) ListProduct(ctx *fiber.Ctx, req dto.ListProductRequest) ([]model.Product, error) {
	products, err := p.repo.NewProductRepo().ListProduct(ctx, req)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (p *productUseCase) CreateProduct(ctx *fiber.Ctx, data *model.Product) (*model.Product, error) {
	data.CreatedAt = time.Now()
	err := p.repo.NewProductRepo().CreateProduct(ctx, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (p *productUseCase) UpdateProduct(ctx *fiber.Ctx, data *model.Product) (*model.Product, error) {
	data.ModifiedAt = time.Now()
	err := p.repo.NewProductRepo().UpdateProductByID(ctx, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (p *productUseCase) CountProduct(ctx *fiber.Ctx) (int32, error) {
	value, err := p.repo.NewProductRepo().CountProduct(ctx)
	if err != nil {
		return 0, err
	}

	return value, nil
}
