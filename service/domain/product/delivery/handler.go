package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/domain/product/usecase"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
)

type ProductHandler struct {
	productUseCase usecase.IProductUseCase
}

func NewProductHandler(productUseCase usecase.IProductUseCase) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
	}
}

func (p *ProductHandler) GetProduct(ctx *fiber.Ctx) error {
	req := new(dto.GetProductByIDRequest)
	if err := ctx.QueryParser(req); err != nil {
		return err
	}

	customer, err := p.productUseCase.GetProductByID(ctx, req.ProductID)
	if err != nil {
		return err
	}

	return ctx.JSON(customer)
}

func (p *ProductHandler) ListProduct(ctx *fiber.Ctx) error {
	req := new(dto.ListProductRequest)
	if err := ctx.QueryParser(req); err != nil {
		return err
	}

	products, err := p.productUseCase.ListProduct(ctx, req.Limit, req.Offset)
	if err != nil {
		return err
	}

	return ctx.JSON(products)
}

func (p *ProductHandler) DeleteProduct(ctx *fiber.Ctx) error {
	req := new(dto.GetProductByIDRequest)
	if err := ctx.QueryParser(req); err != nil {
		return err
	}

	err := p.productUseCase.DeleteProductByID(ctx, req.ProductID)
	if err != nil {
		return err
	}

	return ctx.JSON(dto.DefaultResponse{
		StatusCode: fiber.StatusOK,
	})
}

func (p *ProductHandler) UpdateProduct(ctx *fiber.Ctx) error {
	req := new(model.Product)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	product, err := p.productUseCase.UpdateProduct(ctx, req)
	if err != nil {
		return err
	}

	return ctx.JSON(product)
}

func (p *ProductHandler) CreateProduct(ctx *fiber.Ctx) error {
	req := new(model.Product)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	product, err := p.productUseCase.CreateProduct(ctx, req)
	if err != nil {
		return err
	}

	return ctx.JSON(product)
}
