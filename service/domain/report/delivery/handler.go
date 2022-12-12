package delivery

import "github.com/trungnghia250/malo-api/service/domain/product/usecase"

type ReportHandler struct {
	reportUseCase usecase.IReportUseCase
}

func NewReportHandler(reportUseCase usecase.IReportUseCase) *ReportHandler {
	return &ReportHandler{
		reportUseCase: reportUseCase,
	}
}
