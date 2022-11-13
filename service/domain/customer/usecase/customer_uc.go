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

type ICustomerUseCase interface {
	GetCustomerByID(ctx *fiber.Ctx, customerID string) (*model.Customer, error)
	DeleteCustomersByID(ctx *fiber.Ctx, customerIDs []string) error
	ListCustomer(ctx *fiber.Ctx, req dto.ListCustomerRequest) ([]model.Customer, error)
	CreateCustomer(ctx *fiber.Ctx, data *dto.Customer) (*model.Customer, error)
	UpdateCustomer(ctx *fiber.Ctx, data *dto.Customer) (*model.Customer, error)
	CountCustomer(ctx *fiber.Ctx) (int32, error)
	UpdateCustomerTags(ctx *fiber.Ctx, data dto.UpdateListCustomerRequest) error
	ExportCustomer(ctx *fiber.Ctx, req dto.ExportCustomerRequest) (string, error)
}

type customerUseCase struct {
	repo repo.IRepo
}

func NewCustomerUseCase(repo repo.IRepo) ICustomerUseCase {
	return &customerUseCase{
		repo: repo,
	}
}

func (c *customerUseCase) GetCustomerByID(ctx *fiber.Ctx, customerID string) (*model.Customer, error) {
	customer, err := c.repo.NewCustomerRepo().GetCustomerByID(ctx, customerID)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (c *customerUseCase) DeleteCustomersByID(ctx *fiber.Ctx, customerIDs []string) error {
	err := c.repo.NewCustomerRepo().DeleteCustomersByID(ctx, customerIDs)

	return err
}

func (c *customerUseCase) ListCustomer(ctx *fiber.Ctx, req dto.ListCustomerRequest) ([]model.Customer, error) {
	customers, err := c.repo.NewCustomerRepo().ListCustomer(ctx, req)
	if err != nil {
		return nil, err
	}
	return customers, nil
}

func (c *customerUseCase) CreateCustomer(ctx *fiber.Ctx, data *dto.Customer) (*model.Customer, error) {
	data.CreatedAt = time.Now()

	err := c.repo.NewCustomerRepo().CreateCustomer(ctx, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *customerUseCase) UpdateCustomer(ctx *fiber.Ctx, data *dto.Customer) (*model.Customer, error) {
	data.ModifiedAt = time.Now()

	err := c.repo.NewCustomerRepo().UpdateCustomerByID(ctx, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *customerUseCase) CountCustomer(ctx *fiber.Ctx) (int32, error) {
	value, err := c.repo.NewCustomerRepo().CountCustomer(ctx)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func (c *customerUseCase) UpdateCustomerTags(ctx *fiber.Ctx, data dto.UpdateListCustomerRequest) error {
	err := c.repo.NewCustomerRepo().UpdateListCustomers(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (c *customerUseCase) ExportCustomer(ctx *fiber.Ctx, req dto.ExportCustomerRequest) (string, error) {
	customers, err := c.repo.NewCustomerRepo().ListCustomer(ctx, dto.ListCustomerRequest{
		CustomerIDs: req.CustomerIDs,
		Limit:       int32(len(req.CustomerIDs)),
	})
	if err != nil {
		return "", err
	}

	f := excelize.NewFile()
	categories := map[string]string{
		"A1": "STT",
		"B1": "Mã khách hàng",
		"C1": "Tên khách hàng",
		"D1": "Giới tính",
		"E1": "Số điện thoại",
		"F1": "Email",
		"G1": "Địa chỉ",
		"H1": "Ngày sinh",
		"I1": "Phân loại",
		"J1": "Nguồn khách hàng",
		"K1": "Nhãn",
		"L1": "Ngày tạo",
		"M1": "Ghi chú",
	}
	for k, v := range categories {
		f.SetCellValue("Sheet1", k, v)
	}
	index := 1
	for number, customer := range customers {
		index++
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", index), number+1)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", index), customer.CustomerID)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", index), customer.CustomerName)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", index), customer.Gender)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", index), customer.PhoneNumber)
		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", index), customer.Email)
		f.SetCellValue("Sheet1", fmt.Sprintf("G%d", index), customer.Address)
		f.SetCellValue("Sheet1", fmt.Sprintf("H%d", index), customer.DateOfBirth)
		f.SetCellValue("Sheet1", fmt.Sprintf("I%d", index), customer.CustomerType)
		f.SetCellValue("Sheet1", fmt.Sprintf("J%d", index), customer.CustomerSource)
		f.SetCellValue("Sheet1", fmt.Sprintf("K%d", index), customer.Tags)
		f.SetCellValue("Sheet1", fmt.Sprintf("L%d", index), customer.CreatedAt)
		f.SetCellValue("Sheet1", fmt.Sprintf("M%d", index), customer.Note)
	}
	f.SetActiveSheet(0)
	if err := f.SaveAs("customer.xlsx"); err != nil {
		return "", err
	}
	defer f.Close()

	return f.Path, nil
}
