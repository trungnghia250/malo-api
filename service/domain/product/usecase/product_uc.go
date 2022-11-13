package usecase

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"github.com/trungnghia250/malo-api/service/repo"
	"github.com/xuri/excelize/v2"
	"time"
)

type IProductUseCase interface {
	GetProductByID(ctx *fiber.Ctx, productID string) (*model.Product, error)
	DeleteProductByID(ctx *fiber.Ctx, productIDs []string) error
	ListProduct(ctx *fiber.Ctx, req dto.ListProductRequest) ([]model.Product, error)
	CreateProduct(ctx *fiber.Ctx, data *model.Product) (*model.Product, error)
	UpdateProduct(ctx *fiber.Ctx, data *model.Product) (*model.Product, error)
	CountProduct(ctx *fiber.Ctx) (int32, error)
	ExportProduct(ctx *fiber.Ctx, req dto.ExportProductRequest) (string, error)
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

func (p *productUseCase) DeleteProductByID(ctx *fiber.Ctx, productIDs []string) error {
	err := p.repo.NewProductRepo().DeleteProductByID(ctx, productIDs)

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

func (p *productUseCase) ExportProduct(ctx *fiber.Ctx, req dto.ExportProductRequest) (string, error) {
	products, err := p.repo.NewProductRepo().ListProduct(ctx, dto.ListProductRequest{
		ProductIDs: req.ProductIDs,
		Limit:      int32(len(req.ProductIDs)),
	})
	if err != nil {
		return "", err
	}

	f := excelize.NewFile()
	categories := map[string]string{
		"A1": "STT",
		"B1": "Mã sản phẩm",
		"C1": "Mã SKU",
		"D1": "Tên sản phẩm",
		"E1": "Hình ảnh",
		"F1": "Phân loại",
		"G1": "Mô tả",
		"H1": "Ngày tạo",
		"I1": "Ghi chú",
	}
	for k, v := range categories {
		f.SetCellValue("Sheet1", k, v)
	}
	index := 1
	for number, product := range products {
		index++
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", index), number+1)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", index), product.ProductID)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", index), product.SKU)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", index), product.ProductName)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", index), product.Image)
		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", index), product.Category)
		f.SetCellValue("Sheet1", fmt.Sprintf("G%d", index), product.Description)
		f.SetCellValue("Sheet1", fmt.Sprintf("H%d", index), product.CreatedAt)
		f.SetCellValue("Sheet1", fmt.Sprintf("I%d", index), product.Note)
	}
	f.SetActiveSheet(0)
	if err := f.SaveAs("product.xlsx"); err != nil {
		return "", err
	}
	defer f.Close()

	return f.Path, nil
}
