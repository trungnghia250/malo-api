package usecase

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"github.com/trungnghia250/malo-api/service/repo"
	"github.com/trungnghia250/malo-api/utils"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"time"
)

type IOrderUseCase interface {
	GetOrderByID(ctx *fiber.Ctx, orderID string) (*model.Order, error)
	DeleteOrderByID(ctx *fiber.Ctx, orderID string) error
	ListOrder(ctx *fiber.Ctx, req dto.ListOrderRequest) ([]model.Order, error)
	CreateOrder(ctx *fiber.Ctx, data *dto.Order) (*model.Order, error)
	UpdateOrder(ctx *fiber.Ctx, data *dto.Order) (*model.Order, error)
	CountOrder(ctx *fiber.Ctx) (int32, error)
	ImportOrder(ctx *fiber.Ctx, req dto.ImportOrderRequest) (dto.ImportOrderResponse, error)
}

type orderUseCase struct {
	repo repo.IRepo
}

func NewOrderUseCase(repo repo.IRepo) IOrderUseCase {
	return &orderUseCase{
		repo: repo,
	}
}

func (o *orderUseCase) GetOrderByID(ctx *fiber.Ctx, orderID string) (*model.Order, error) {
	order, err := o.repo.NewOrderRepo().GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (o *orderUseCase) DeleteOrderByID(ctx *fiber.Ctx, orderID string) error {
	err := o.repo.NewOrderRepo().DeleteOrderByID(ctx, orderID)

	return err
}

func (o *orderUseCase) ListOrder(ctx *fiber.Ctx, req dto.ListOrderRequest) ([]model.Order, error) {
	orders, err := o.repo.NewOrderRepo().ListOrder(ctx, req)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (o *orderUseCase) CreateOrder(ctx *fiber.Ctx, data *dto.Order) (*model.Order, error) {
	_, err := o.repo.NewOrderRepo().CreateOrder(ctx, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (o *orderUseCase) UpdateOrder(ctx *fiber.Ctx, data *dto.Order) (*model.Order, error) {
	data.ModifiedAt = time.Now()
	err := o.repo.NewOrderRepo().UpdateOrderByID(ctx, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (o *orderUseCase) CountOrder(ctx *fiber.Ctx) (int32, error) {
	value, err := o.repo.NewOrderRepo().CountOrder(ctx)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func (o *orderUseCase) ImportOrder(ctx *fiber.Ctx, req dto.ImportOrderRequest) (dto.ImportOrderResponse, error) {
	file, err := req.File.Open()
	if err != nil {
		return dto.ImportOrderResponse{}, err
	}
	defer file.Close()

	excelFile, err := excelize.OpenReader(file)

	rows, err := excelFile.GetRows("order")
	if err != nil {
		return dto.ImportOrderResponse{}, err
	}
	var listOrder []dto.Order
	var tempOrder dto.Order

	for _, row := range rows[1:] {
		if row[0] == "" {
			tempOrder.Items = append(tempOrder.Items, dto.Item{
				ProductName:    row[9],
				SKU:            row[8],
				Quantity:       convertStringToInt32(row[10]),
				UnitPrice:      convertStringToInt32(row[11]),
				TotalDiscount:  convertStringToInt32(row[12]),
				TotalTaxAmount: convertStringToInt32(row[13]),
				Subtotal:       convertStringToInt32(row[14]),
			})
			continue
		}
		if tempOrder.OrderID != "" {
			listOrder = append(listOrder, tempOrder)
		}

		thisOrder := dto.Order{
			OrderID:      row[1],
			CustomerName: row[3],
			PhoneNumber:  row[5],
			Email:        row[6],
			Address:      row[7],
			Source:       row[2],
			Status:       row[21],
			Items: []dto.Item{
				{
					ProductName:    row[9],
					SKU:            row[8],
					Quantity:       convertStringToInt32(row[10]),
					UnitPrice:      convertStringToInt32(row[11]),
					TotalDiscount:  convertStringToInt32(row[12]),
					TotalTaxAmount: convertStringToInt32(row[13]),
					Subtotal:       convertStringToInt32(row[14]),
				},
			},
			VoucherCode:          row[15],
			TotalLineItemsAmount: convertStringToInt32(row[16]),
			ShippingPrice:        convertStringToInt32(row[17]),
			TotalDiscount:        convertStringToInt32(row[18]),
			TotalTaxAmount:       convertStringToInt32(row[19]),
			TotalOrderAmount:     convertStringToInt32(row[20]),
			Note:                 dataEndLine(row),
			CreateAt:             time.Now(),
			ModifiedAt:           time.Now(),
		}
		tempOrder = thisOrder
	}
	listOrder = append(listOrder, tempOrder)

	resp, err := o.CheckOrderImport(ctx, listOrder, req.Action, req.CheckDupCol)
	if err != nil {
		return dto.ImportOrderResponse{}, err
	}
	return resp, nil
}

func (o *orderUseCase) CheckOrderImport(ctx *fiber.Ctx, data []dto.Order, action string, checkDupCol []string) (dto.ImportOrderResponse, error) {
	totalInsert := int32(0)
	totalUpdate := int32(0)
	totalIgnore := int32(0)
	var orderIDs []string
	for _, order := range data {
		matching := bson.M{}
		if utils.IsStringContains(checkDupCol, "source") {
			matching["source"] = order.Source
		}
		if utils.IsStringContains(checkDupCol, "order_id") {
			matching["order_id"] = order.OrderID
		}

		if action == "override" {
			insert, update, _ := o.repo.NewOrderRepo().UpsertOrder(ctx, matching, order)
			totalInsert += insert
			totalUpdate += update
			if !utils.IsStringContains(orderIDs, order.OrderID) {
				orderIDs = append(orderIDs, order.OrderID)
			}
		}

		if action == "ignore" {
			isExist, _ := o.repo.NewOrderRepo().CheckOrderExist(ctx, matching)
			if !isExist {
				insert, _ := o.repo.NewOrderRepo().CreateOrder(ctx, &order)
				totalInsert += insert
				orderIDs = append(orderIDs, order.OrderID)
			} else {
				totalIgnore += 1
			}
		}

	}

	orders, err := o.repo.NewOrderRepo().ListOrder(ctx, dto.ListOrderRequest{
		OrderIDs: orderIDs,
		Limit:    100,
	})
	if err != nil {
		return dto.ImportOrderResponse{}, err
	}
	resp := dto.ImportOrderResponse{
		Scan:    int32(len(data)),
		Success: totalInsert + totalUpdate,
		Insert:  totalInsert,
		Update:  totalUpdate,
		Ignore:  totalIgnore,
		Data:    orders,
	}
	return resp, nil
}

func convertStringToInt32(data string) int32 {
	value, _ := strconv.Atoi(data)

	return int32(value)
}

func dataEndLine(data []string) string {
	if len(data) < 23 {
		return ""
	}
	return data[22]
}
