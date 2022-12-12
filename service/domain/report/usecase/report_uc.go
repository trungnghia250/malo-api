package usecase

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"github.com/trungnghia250/malo-api/service/repo"
	"time"
)

type IReportUseCase interface {
	GetReportByCategory(ctx *fiber.Ctx, req dto.GetReportRequest) (interface{}, error)
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

	switch req.Type {
	case "customer":
		reports, err := r.repo.NewCustomerReportRepo().GetCustomerReport(ctx, start, end)
		if err != nil {
			return nil, err
		}
		var orders, success, process, cancel, revenue int32
		for _, report := range reports {
			orders += report.TotalOrders
			success += report.SuccessOrders
			process += report.ProcessingOrders
			cancel += report.CancelOrders
			revenue += report.TotalRevenue
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
			},
		}, nil
	case "product":
		reports, err := r.repo.NewProductReportRepo().GetProductReport(ctx, start, end)
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

func ListCustomerReportsPaginate(records []dto.CustomerReport, limit, offset int32) []dto.CustomerReport {
	if offset > int32(len(records)) {
		offset = int32(len(records))
	}
	end := offset + limit
	if end > int32(len(records)) {
		end = int32(len(records))
	}
	return records[offset:limit]
}

func ListProductReportsPaginate(records []dto.ProductReport, limit, offset int32) []dto.ProductReport {
	if offset > int32(len(records)) {
		offset = int32(len(records))
	}
	end := offset + limit
	if end > int32(len(records)) {
		end = int32(len(records))
	}
	return records[offset:limit]
}
