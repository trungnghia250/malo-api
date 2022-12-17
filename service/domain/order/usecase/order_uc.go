package usecase

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"github.com/trungnghia250/malo-api/service/model/transform"
	"github.com/trungnghia250/malo-api/service/repo"
	"github.com/trungnghia250/malo-api/utils"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strconv"
	"time"
)

type IOrderUseCase interface {
	GetOrderByID(ctx *fiber.Ctx, orderID string) (*model.Order, error)
	DeleteOrderByID(ctx *fiber.Ctx, orderIDs []string) error
	ListOrder(ctx *fiber.Ctx, req dto.ListOrderRequest) ([]model.Order, error)
	CreateOrder(ctx *fiber.Ctx, data *dto.Order) (*model.Order, error)
	UpdateOrder(ctx *fiber.Ctx, data *dto.Order) (*model.Order, error)
	CountOrder(ctx *fiber.Ctx) (int32, error)
	ImportOrder(ctx *fiber.Ctx, req dto.ImportOrderRequest) (dto.ImportOrderResponse, error)
	ExportOrder(ctx *fiber.Ctx, req dto.ExportOrderRequest) (string, error)
	SyncOrder(ctx *fiber.Ctx, req dto.SyncOrderRequest) (dto.ImportOrderResponse, error)
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

func (o *orderUseCase) DeleteOrderByID(ctx *fiber.Ctx, orderIDs []string) error {
	err := o.repo.NewOrderRepo().DeleteOrderByID(ctx, orderIDs)

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
	orderID, err := o.repo.NewCounterRepo().GetSequenceNextValue(ctx, "order_id")
	if err != nil {
		return nil, err
	}
	data.OrderID = fmt.Sprintf("O%d", orderID)
	_, err = o.repo.NewOrderRepo().CreateOrder(ctx, data)
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
				ProductName:    row[10],
				SKU:            row[9],
				Quantity:       convertStringToInt32(row[11]),
				UnitPrice:      convertStringToInt32(row[12]),
				TotalDiscount:  convertStringToInt32(row[13]),
				TotalTaxAmount: convertStringToInt32(row[14]),
				Subtotal:       convertStringToInt32(row[15]),
			})
			continue
		}
		if tempOrder.OrderID != "" {
			listOrder = append(listOrder, tempOrder)
		}
		status := ""
		if row[22] == "Đang xử lý" {
			status = "processing"
		}
		if row[22] == "Đã hoàn thành" {
			status = "success"
		}
		if row[22] == "Đã huỷ" {
			status = "cancel"
		}
		thisOrder := dto.Order{
			OrderID:      row[1],
			CustomerName: row[3],
			PhoneNumber:  row[5],
			Email:        row[6],
			Address:      row[7],
			Source:       row[2],
			Status:       status,
			Items: []dto.Item{
				{
					ProductName:    row[10],
					SKU:            row[9],
					Quantity:       convertStringToInt32(row[11]),
					UnitPrice:      convertStringToInt32(row[12]),
					TotalDiscount:  convertStringToInt32(row[13]),
					TotalTaxAmount: convertStringToInt32(row[14]),
					Subtotal:       convertStringToInt32(row[15]),
				},
			},
			VoucherCode:          row[16],
			TotalLineItemsAmount: convertStringToInt32(row[17]),
			ShippingPrice:        convertStringToInt32(row[18]),
			TotalDiscount:        convertStringToInt32(row[19]),
			TotalTaxAmount:       convertStringToInt32(row[20]),
			TotalOrderAmount:     convertStringToInt32(row[21]),
			Note:                 dataEndLine(row),
			CreatedAt:            time.Now(),
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

	if len(orderIDs) == 0 {
		return dto.ImportOrderResponse{
			Scan:    int32(len(data)),
			Success: totalInsert + totalUpdate,
			Insert:  totalInsert,
			Update:  totalUpdate,
			Ignore:  totalIgnore,
			Data:    nil,
		}, nil
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
	if len(data) < 24 {
		return ""
	}
	return data[23]
}

func (o *orderUseCase) ExportOrder(ctx *fiber.Ctx, req dto.ExportOrderRequest) (string, error) {
	query := dto.ListOrderRequest{}

	if len(req.Filter) > 0 {
		json.Unmarshal([]byte(req.Filter), &query)
	}

	if len(req.OrderIDs) > 0 {
		query = dto.ListOrderRequest{
			OrderIDs: req.OrderIDs,
			Limit:    int32(len(req.OrderIDs)),
		}
	}

	orders, err := o.repo.NewOrderRepo().ListOrder(ctx, query)
	if err != nil {
		return "", err
	}

	f := excelize.NewFile()
	categories := map[string]string{
		"A1": "STT",
		"B1": "Mã đơn hàng",
		"C1": "Nguồn đơn hàng",
		"D1": "Tên khách hàng",
		"E1": "Giới tính",
		"F1": "Số điện thoại",
		"G1": "Email",
		"H1": "Địa chỉ",
		"I1": "Mã SKU",
		"J1": "Tên sản phẩm",
		"K1": "Số lượng",
		"L1": "Đơn giá",
		"M1": "Giảm giá",
		"N1": "Thuế",
		"O1": "Thành tiền",
		"P1": "Mã giảm giá",
		"Q1": "Tổng giá trị sản phẩm",
		"R1": "Phí vận chuyển",
		"S1": "Tổng giảm giá",
		"T1": "Tổng thuế",
		"U1": "Giá trị đơn hàng",
		"V1": "Trạng thái",
		"W1": "Ghi Chú",
	}
	for k, v := range categories {
		f.SetCellValue("Sheet1", k, v)
	}
	index := 1
	for number, order := range orders {
		index++
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", index), number+1)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", index), order.OrderID)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", index), order.Source)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", index), order.CustomerName)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", index), order.Gender)
		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", index), order.PhoneNumber)
		f.SetCellValue("Sheet1", fmt.Sprintf("G%d", index), order.Email)
		f.SetCellValue("Sheet1", fmt.Sprintf("H%d", index), order.Address)
		f.SetCellValue("Sheet1", fmt.Sprintf("I%d", index), order.Items[0].SKU)
		f.SetCellValue("Sheet1", fmt.Sprintf("J%d", index), order.Items[0].ProductName)
		f.SetCellValue("Sheet1", fmt.Sprintf("K%d", index), order.Items[0].Quantity)
		f.SetCellValue("Sheet1", fmt.Sprintf("L%d", index), order.Items[0].UnitPrice)
		f.SetCellValue("Sheet1", fmt.Sprintf("M%d", index), order.Items[0].TotalDiscount)
		f.SetCellValue("Sheet1", fmt.Sprintf("N%d", index), order.Items[0].TotalTaxAmount)
		f.SetCellValue("Sheet1", fmt.Sprintf("O%d", index), order.Items[0].Subtotal)
		f.SetCellValue("Sheet1", fmt.Sprintf("P%d", index), order.VoucherCode)
		f.SetCellValue("Sheet1", fmt.Sprintf("Q%d", index), order.TotalLineItemsAmount)
		f.SetCellValue("Sheet1", fmt.Sprintf("R%d", index), order.ShippingPrice)
		f.SetCellValue("Sheet1", fmt.Sprintf("S%d", index), order.TotalDiscount)
		f.SetCellValue("Sheet1", fmt.Sprintf("T%d", index), order.TotalTaxAmount)
		f.SetCellValue("Sheet1", fmt.Sprintf("U%d", index), order.TotalOrderAmount)
		f.SetCellValue("Sheet1", fmt.Sprintf("V%d", index), order.Status)
		f.SetCellValue("Sheet1", fmt.Sprintf("W%d", index), order.Note)

		if len(order.Items) > 1 {
			for _, item := range order.Items[1:] {
				index++
				f.SetCellValue("Sheet1", fmt.Sprintf("I%d", index), item.SKU)
				f.SetCellValue("Sheet1", fmt.Sprintf("J%d", index), item.ProductName)
				f.SetCellValue("Sheet1", fmt.Sprintf("K%d", index), item.Quantity)
				f.SetCellValue("Sheet1", fmt.Sprintf("L%d", index), item.UnitPrice)
				f.SetCellValue("Sheet1", fmt.Sprintf("M%d", index), item.TotalDiscount)
				f.SetCellValue("Sheet1", fmt.Sprintf("N%d", index), item.TotalTaxAmount)
				f.SetCellValue("Sheet1", fmt.Sprintf("O%d", index), item.Subtotal)
			}
		}
	}

	f.SetActiveSheet(0)
	if err := f.SaveAs("order.xlsx"); err != nil {
		return "", err
	}
	defer f.Close()

	return f.Path, nil
}

func (o *orderUseCase) SyncOrder(ctx *fiber.Ctx, req dto.SyncOrderRequest) (dto.ImportOrderResponse, error) {
	switch req.Source {
	case "SAPO":
		client := &http.Client{}
		request, err := http.NewRequest("GET", "https://malo25.mysapo.net/admin/orders.json", nil)
		if err != nil {
			return dto.ImportOrderResponse{}, err
		}
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("X-Sapo-Access-Token", "007eba624916497080f1c8ae3e3d03da")

		res, err := client.Do(request)
		if err != nil {
			return dto.ImportOrderResponse{}, err
		}
		defer func() {
			err := res.Body.Close()
			if err != nil {
				fmt.Printf("failed to closed connection: %v", err)
			}
		}()
		var result dto.SapoOrdersResponse
		err = json.NewDecoder(res.Body).Decode(&result)
		if err != nil {
			return dto.ImportOrderResponse{}, err
		}
		maloOrders := transform.ConvertSapoOrder(result.Orders)

		resp, err := o.CheckOrderImport(ctx, maloOrders, req.Action, req.ArrayCol)
		if err != nil {
			return dto.ImportOrderResponse{}, err
		}
		return resp, nil

	default:
		return dto.ImportOrderResponse{}, nil
	}
}
