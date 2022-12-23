package usecase

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"github.com/trungnghia250/malo-api/service/repo"
	"github.com/xuri/excelize/v2"
	"strings"
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
	ImportCustomer(ctx *fiber.Ctx, req dto.ImportCustomerRequest) (dto.ImportCustomerResponse, error)
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
	customerID, err := c.repo.NewCounterRepo().GetSequenceNextValue(ctx, "customer_id")
	if err != nil {
		return nil, err
	}
	data.CustomerID = fmt.Sprintf("C%d", customerID)
	err = c.repo.NewCustomerRepo().CreateCustomer(ctx, data)
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
	query := dto.ListCustomerRequest{}

	if len(req.Filter) > 0 {
		json.Unmarshal([]byte(req.Filter), &query)
	}

	if len(req.CustomerIDs) > 0 {
		query = dto.ListCustomerRequest{
			CustomerIDs: req.CustomerIDs,
			Limit:       int32(len(req.CustomerIDs)),
		}
	}

	customers, err := c.repo.NewCustomerRepo().ListCustomer(ctx, query)
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

func (c *customerUseCase) ImportCustomer(ctx *fiber.Ctx, req dto.ImportCustomerRequest) (dto.ImportCustomerResponse, error) {
	file, err := req.File.Open()
	if err != nil {
		return dto.ImportCustomerResponse{}, err
	}
	defer file.Close()

	excelFile, err := excelize.OpenReader(file)

	rows, err := excelFile.GetRows("customer")
	if err != nil {
		return dto.ImportCustomerResponse{}, err
	}
	var listCustomer []dto.Customer

	for _, row := range rows[1:] {
		dob, _ := time.Parse("2006/01/02", row[7])
		var gender string
		if strings.ToLower(row[2]) == "nam" {
			gender = "male"
		} else {
			gender = "female"
		}

		thisCustomer := dto.Customer{
			CustomerName:   row[1],
			Gender:         gender,
			PhoneNumber:    row[3],
			Email:          row[4],
			Address:        row[5],
			Province:       row[6],
			DateOfBirth:    dob,
			CustomerSource: row[9],
			Note:           dataEndLine(row),
			CreatedAt:      time.Now(),
		}
		var tags []string
		if len(row[8]) > 0 {
			tags = strings.Split(row[8], ",")
			thisCustomer.Tags = tags
		}
		listCustomer = append(listCustomer, thisCustomer)

	}

	resp, err := c.CheckCustomerImport(ctx, listCustomer)
	if err != nil {
		return dto.ImportCustomerResponse{}, err
	}
	return resp, nil
}

func dataEndLine(data []string) string {
	if len(data) < 11 {
		return ""
	}
	return data[10]
}

func (c *customerUseCase) CheckCustomerImport(ctx *fiber.Ctx, data []dto.Customer) (dto.ImportCustomerResponse, error) {
	totalInsert := int32(0)
	totalUpdate := int32(0)
	totalIgnore := int32(0)
	var customerIDs []string
	for _, customer := range data {
		cus, _ := c.repo.NewCustomerRepo().GetCustomerByPhone(ctx, customer.PhoneNumber)
		if cus.PhoneNumber == "" {
			customerID, _ := c.repo.NewCounterRepo().GetSequenceNextValue(ctx, "customer_id")
			customer.CustomerID = fmt.Sprintf("C%d", customerID)
			_ = c.repo.NewCustomerRepo().CreateCustomer(ctx, &customer)
			customerIDs = append(customerIDs, customer.CustomerID)
			totalInsert += 1
			continue
		}
		customerIDs = append(customerIDs, cus.CustomerID)
		_ = c.repo.NewCustomerRepo().UpdateCustomerByID(ctx, &customer)
		totalUpdate += 1
	}

	if len(customerIDs) == 0 {
		return dto.ImportCustomerResponse{
			Scan:    int32(len(data)),
			Success: totalInsert + totalUpdate,
			Insert:  totalInsert,
			Update:  totalUpdate,
			Ignore:  totalIgnore,
			Data:    nil,
		}, nil
	}

	customers, err := c.repo.NewCustomerRepo().ListCustomer(ctx, dto.ListCustomerRequest{
		CustomerIDs: customerIDs,
		Limit:       100,
	})
	if err != nil {
		return dto.ImportCustomerResponse{}, err
	}
	resp := dto.ImportCustomerResponse{
		Scan:    int32(len(data)),
		Success: totalInsert + totalUpdate,
		Insert:  totalInsert,
		Update:  totalUpdate,
		Ignore:  totalIgnore,
		Data:    customers,
	}
	return resp, nil
}
