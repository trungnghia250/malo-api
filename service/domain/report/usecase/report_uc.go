package usecase

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"github.com/trungnghia250/malo-api/service/repo"
	"github.com/trungnghia250/malo-api/utils"
	"github.com/xuri/excelize/v2"
	"time"
)

type IReportUseCase interface {
	GetReportByCategory(ctx *fiber.Ctx, req dto.GetReportRequest) (interface{}, error)
	GetDashboard(ctx *fiber.Ctx, req dto.GetDashBoardRequest) (interface{}, error)
}

type reportUseCase struct {
	repo repo.IRepo
}

func NewReportUseCase(repo repo.IRepo) IReportUseCase {
	return &reportUseCase{
		repo: repo,
	}
}

func (r *reportUseCase) GetReportByCategory(ctx *fiber.Ctx, req dto.GetReportRequest) (interface{}, error) {
	end, _ := time.Parse("02/01/2006", req.EndTime)
	start, _ := time.Parse("02/01/2006", req.StartTime)

	if len(req.GroupID) > 0 {
		group, _ := r.repo.NewCustomerGroupRepo().GetCustomerGroupByID(ctx, req.GroupID)
		customers, _ := r.repo.NewCustomerRepo().ListCustomer(ctx, dto.ListCustomerRequest{
			Limit:       int32(len(group.CustomerIDs)),
			Offset:      0,
			CustomerIDs: group.CustomerIDs,
		})
		var phones []string
		for _, customer := range customers {
			phones = append(phones, customer.PhoneNumber)
		}
		req.Phone = phones
	}

	switch req.Type {
	case "customer":
		reports, err := r.repo.NewCustomerReportRepo().GetCustomerReport(ctx, start, end, req)
		if err != nil {
			return nil, err
		}

		var orders, success, process, cancel, revenue, isNew int32
		for _, report := range reports {
			orders += report.TotalOrders
			success += report.SuccessOrders
			process += report.ProcessingOrders
			cancel += report.CancelOrders
			revenue += report.TotalRevenue
			isNew += report.New
		}
		if req.Export {
			f := excelize.NewFile()
			categories := map[string]string{
				"A1": "STT",
				"B1": "Tên khách hàng",
				"C1": "Số điện thoại",
				"D1": "Email",
				"E1": "Tổng đơn hàng",
				"F1": "Đơn hàng thành công",
				"G1": "Đơn hàng đang xử lý",
				"H1": "Đơn hàng bị huỷ",
				"I1": "Tổng doanh thu",
			}
			for k, v := range categories {
				f.SetCellValue("Sheet1", k, v)
			}
			index := 1
			for number, report := range reports {
				index++
				f.SetCellValue("Sheet1", fmt.Sprintf("A%d", index), number+1)
				f.SetCellValue("Sheet1", fmt.Sprintf("B%d", index), report.Name)
				f.SetCellValue("Sheet1", fmt.Sprintf("C%d", index), report.Phone)
				f.SetCellValue("Sheet1", fmt.Sprintf("D%d", index), report.Email)
				f.SetCellValue("Sheet1", fmt.Sprintf("E%d", index), report.TotalOrders)
				f.SetCellValue("Sheet1", fmt.Sprintf("F%d", index), report.SuccessOrders)
				f.SetCellValue("Sheet1", fmt.Sprintf("G%d", index), report.ProcessingOrders)
				f.SetCellValue("Sheet1", fmt.Sprintf("H%d", index), report.CancelOrders)
				f.SetCellValue("Sheet1", fmt.Sprintf("I%d", index), report.TotalRevenue)
			}
			index++
			f.SetCellValue("Sheet1", fmt.Sprintf("A%d", index), "Tổng cộng")
			totalLine := map[string]int32{
				fmt.Sprintf("E%d", index): orders,
				fmt.Sprintf("F%d", index): success,
				fmt.Sprintf("G%d", index): process,
				fmt.Sprintf("H%d", index): cancel,
				fmt.Sprintf("I%d", index): revenue,
			}
			for k, v := range totalLine {
				f.SetCellValue("Sheet1", k, v)
			}
			f.SetActiveSheet(0)
			if err = f.SaveAs("customer_report.xlsx"); err != nil {
				return "", err
			}
			defer f.Close()

			return f.Path, nil
		}
		return dto.CustomerReportResponse{
			Data:  ListCustomerReportsPaginate(reports, req.Limit, req.Offset),
			Count: int32(len(reports)),
			Total: dto.CustomerReport{
				Name:             "Tổng cộng",
				TotalOrders:      orders,
				SuccessOrders:    success,
				ProcessingOrders: process,
				CancelOrders:     cancel,
				TotalRevenue:     revenue,
				New:              isNew,
				Return:           ComebackCustomer(reports) - isNew,
			},
		}, nil
	case "product":
		reports, err := r.repo.NewProductReportRepo().GetProductReport(ctx, start, end, req)
		if err != nil {
			return nil, err
		}
		var orders, success, process, cancel, revenue, sales int32
		for _, report := range reports {
			orders += report.TotalOrders
			success += report.SuccessOrders
			process += report.ProcessingOrders
			cancel += report.CancelOrders
			revenue += report.TotalRevenue
			sales += report.TotalSales
		}
		if req.Export {
			f := excelize.NewFile()
			categories := map[string]string{
				"A1": "STT",
				"B1": "Mã SKU",
				"C1": "Tên sản phẩm",
				"D1": "Số lượng bán ra",
				"E1": "Tổng đơn hàng",
				"F1": "Đơn hàng thành công",
				"G1": "Đơn hàng đang xử lý",
				"H1": "Đơn hàng bị huỷ",
				"I1": "Tổng doanh thu",
			}
			for k, v := range categories {
				f.SetCellValue("Sheet1", k, v)
			}
			index := 1
			for number, report := range reports {
				index++
				f.SetCellValue("Sheet1", fmt.Sprintf("A%d", index), number+1)
				f.SetCellValue("Sheet1", fmt.Sprintf("B%d", index), report.SKU)
				f.SetCellValue("Sheet1", fmt.Sprintf("C%d", index), report.Name)
				f.SetCellValue("Sheet1", fmt.Sprintf("D%d", index), report.TotalSales)
				f.SetCellValue("Sheet1", fmt.Sprintf("E%d", index), report.TotalOrders)
				f.SetCellValue("Sheet1", fmt.Sprintf("F%d", index), report.SuccessOrders)
				f.SetCellValue("Sheet1", fmt.Sprintf("G%d", index), report.ProcessingOrders)
				f.SetCellValue("Sheet1", fmt.Sprintf("H%d", index), report.CancelOrders)
				f.SetCellValue("Sheet1", fmt.Sprintf("I%d", index), report.TotalRevenue)
			}
			index++
			f.SetCellValue("Sheet1", fmt.Sprintf("A%d", index), "Tổng cộng")
			totalLine := map[string]int32{
				fmt.Sprintf("D%d", index): sales,
				fmt.Sprintf("E%d", index): orders,
				fmt.Sprintf("F%d", index): success,
				fmt.Sprintf("G%d", index): process,
				fmt.Sprintf("H%d", index): cancel,
				fmt.Sprintf("I%d", index): revenue,
			}
			for k, v := range totalLine {
				f.SetCellValue("Sheet1", k, v)
			}
			f.SetActiveSheet(0)
			if err = f.SaveAs("product_report.xlsx"); err != nil {
				return "", err
			}
			defer f.Close()

			return f.Path, nil
		}
		return dto.ProductReportResponse{
			Data:  ListProductReportsPaginate(reports, req.Limit, req.Offset),
			Count: int32(len(reports)),
			Total: dto.ProductReport{
				Name:             "Tổng cộng",
				TotalOrders:      orders,
				SuccessOrders:    success,
				ProcessingOrders: process,
				CancelOrders:     cancel,
				TotalRevenue:     revenue,
				TotalSales:       sales,
			},
		}, nil
	}

	return nil, nil
}

func (r *reportUseCase) GetDashboard(ctx *fiber.Ctx, req dto.GetDashBoardRequest) (interface{}, error) {
	end, _ := time.Parse("02/01/2006", req.EndTime)
	start, _ := time.Parse("02/01/2006", req.StartTime)
	var phones []string

	if len(req.GroupID) > 0 {
		groups, _ := r.repo.NewCustomerGroupRepo().ListCustomerGroup(ctx, dto.ListCustomerGroupRequest{
			Limit:  int32(len(req.GroupID)),
			Offset: 0,
			IDs:    req.GroupID,
		})

		for _, group := range groups {
			req.CustomerID = append(req.CustomerID, group.CustomerIDs...)
		}
	}

	if len(req.CustomerID) > 0 {
		customers, _ := r.repo.NewCustomerRepo().ListCustomer(ctx, dto.ListCustomerRequest{
			Limit:       int32(len(req.CustomerID)),
			Offset:      0,
			CustomerIDs: req.CustomerID,
		})
		for _, customer := range customers {
			phones = append(phones, customer.PhoneNumber)
		}
	}

	if req.Type == "product" {
		reports, err := r.repo.NewProductReportRepo().GetDashboard(ctx, start, end, phones)
		if err != nil {
			return nil, err
		}
		var orders, revenue int32
		for _, report := range reports {
			orders += report.TotalOrders
			revenue += report.TotalRevenue
		}
		return dto.CustomerReportResponse{
			Data:  reports,
			Count: int32(len(reports)),
			Total: dto.CustomerReport{
				Name:         "Tổng cộng",
				TotalOrders:  orders,
				TotalRevenue: revenue,
			},
		}, nil
	}

	reports, err := r.repo.NewCustomerReportRepo().GetDashboard(ctx, start, end, phones)
	if err != nil {
		return nil, err
	}

	var orders, revenue, isNew, comeback int32
	for i, report := range reports {
		orders += report.TotalOrders
		//success += report.SuccessOrders
		//process += report.ProcessingOrders
		//cancel += report.CancelOrders
		revenue += report.TotalRevenue
		isNew += report.New
		reports[i].Return = report.TotalOrders - report.New
		comeback += reports[i].Return
	}
	totalReturn, _ := r.repo.NewCustomerReportRepo().GetReturn(ctx, start, end)
	return dto.CustomerReportResponse{
		Data:  reports,
		Count: int32(len(reports)),
		Total: dto.CustomerReport{
			Name:         "Tổng cộng",
			TotalOrders:  orders,
			TotalRevenue: revenue,
			New:          isNew,
			Return:       totalReturn - isNew,
		},
	}, nil

}

func ListCustomerReportsPaginate(records []dto.CustomerReport, limit, offset int32) []dto.CustomerReport {
	if limit == 0 && offset == 0 {
		return records
	}
	if offset > int32(len(records)) {
		offset = int32(len(records))
	}
	end := offset + limit
	if end > int32(len(records)) {
		end = int32(len(records))
	}
	return records[offset:end]
}

func ListProductReportsPaginate(records []dto.ProductReport, limit, offset int32) []dto.ProductReport {
	if limit == 0 && offset == 0 {
		return records
	}
	if offset > int32(len(records)) {
		offset = int32(len(records))
	}
	end := offset + limit
	if end > int32(len(records)) {
		end = int32(len(records))
	}
	return records[offset:end]
}

func ComebackCustomer(data []dto.CustomerReport) int32 {
	var customer []string
	for _, record := range data {
		if !utils.IsStringContains(customer, record.Name) {
			customer = append(customer, record.Name)
		}
	}

	return int32(len(customer))
}
