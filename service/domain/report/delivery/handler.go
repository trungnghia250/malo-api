package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/domain/report/usecase"
	"github.com/trungnghia250/malo-api/service/model/dto"
)

type ReportHandler struct {
	reportUseCase usecase.IReportUseCase
}

func NewReportHandler(reportUseCase usecase.IReportUseCase) *ReportHandler {
	return &ReportHandler{
		reportUseCase: reportUseCase,
	}
}

func (r *ReportHandler) GetReport(ctx *fiber.Ctx) error {
	req := new(dto.GetReportRequest)
	if err := ctx.QueryParser(req); err != nil {
		return err
	}

	result, err := r.reportUseCase.GetReportByCategory(ctx, *req)
	if err != nil {
		return err
	}
	if req.Export {
		return ctx.SendFile(result.(string))
	}
	return ctx.JSON(result)
}
