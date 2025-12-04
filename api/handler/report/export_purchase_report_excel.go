package report

import (
	"log/slog"
	reportCommand "mini-erp-backend/api/service/report/command"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

// ExportPurchaseReportExcel
//
//	@Summary		Export purchase report to Excel
//	@Description	Export purchase order summary by month to Excel file
//	@Tags			Report
//	@Produce		application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
//	@Param			month	query	string	true	"Month (MM-YYYY)"
//	@Success		200	{file}	file
//	@Failure		400	{object}	api.ErrorResponse
//	@Failure		500	{object}	api.ErrorResponse
//	@Router			/reports/purchase-summary/export [get]
func ExportPurchaseReportExcel(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		monthStr := c.Query("month")
		if monthStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "month query parameter is required (format: MM-YYYY)",
			})
		}

		monthDate, err := time.Parse("01-2006", monthStr)
		if err != nil {
			logger.Error("Invalid month format", "month", monthStr, "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid month format (expected: MM-YYYY)",
			})
		}

		req := &reportCommand.ExportPurchaseReportExcelRequest{Year: monthDate.Year(), Month: int(monthDate.Month())}

		result, err := mediatr.Send[*reportCommand.ExportPurchaseReportExcelRequest, *reportCommand.ExportPurchaseReportExcelResult](c.Context(), req)
		if err != nil {
			logger.Error("Failed to export purchase report Excel", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to export purchase report",
			})
		}

		c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Set("Content-Disposition", "attachment; filename="+result.Filename)
		return c.Send(result.Data)
	}
}