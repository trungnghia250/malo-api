package delivery

import "github.com/gofiber/fiber/v2"

func (r *ReportHandler) InternalReportAPIRoute(router fiber.Router) {
	router.Get("/report", r.GetReport)
	router.Get("/dashboard", r.GetDashboard)
}
