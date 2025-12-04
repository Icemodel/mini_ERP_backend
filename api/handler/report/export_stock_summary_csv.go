package report

import (
	"log/slog"
	reportCommand "mini-erp-backend/api/service/report/command"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

// ExportStockSummaryCSV
//
//	@Summary		Export stock summary to CSV
//	@Description	Export current stock summary to CSV file
//	@Tags			Report
//	@Produce		text/csv
//	@Success		200	{file}	file
//	@Failure		500	{object}	api.ErrorResponse
//	@Router			/reports/stock-summary/export [get]
func ExportStockSummaryCSV(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := &reportCommand.ExportStockSummaryCSVRequest{}

		result, err := mediatr.Send[*reportCommand.ExportStockSummaryCSVRequest, *reportCommand.ExportStockSummaryCSVResult](c.Context(), req)
		if err != nil {
			logger.Error("Failed to export stock summary CSV", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to export stock summary",
			})
		}

		c.Set("Content-Type", "text/csv")
		c.Set("Content-Disposition", "attachment; filename="+result.Filename)
		return c.Send(result.Data)
	}
}