package report

import (
	"log/slog"
	reportCommand "mini-erp-backend/api/service/report/command"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

// ExportStockMovementExcel
//
// 	@Summary		Export stock movements to Excel
// 	@Description	Export stock movements within a date range to Excel file
// 	@Tags			Report
// 	@Produce		application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// 	@Param			from	query	string	true	"From date (DD-MM-YYYY)"
// 	@Param			to		query	string	true	"To date (DD-MM-YYYY)"
// 	@Success		200	{file}	file
// 	@Failure		400	{object}	fiber.Map
// 	@Failure		500	{object}	fiber.Map
// 	@Router			/api/v1/reports/stock-movements/export [get]
func ExportStockMovementExcel(logger *slog.Logger) fiber.Handler {
    return func(c *fiber.Ctx) error {
        fromStr := c.Query("from")
        toStr := c.Query("to")

        if fromStr == "" || toStr == "" {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": "from and to query parameters are required (format: DD-MM-YYYY)",
            })
        }

        fromDate, err := time.Parse("02-01-2006", fromStr)
        if err != nil {
            logger.Error("Invalid from date", "from", fromStr, "error", err)
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": "Invalid from date format (expected: DD-MM-YYYY)",
            })
        }
        toDate, err := time.Parse("02-01-2006", toStr)
        if err != nil {
            logger.Error("Invalid to date", "to", toStr, "error", err)
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": "Invalid to date format (expected: DD-MM-YYYY)",
            })
        }
        toDate = toDate.Add(24*time.Hour - time.Second)

        req := &reportCommand.ExportStockMovementExcelRequest{FromDate: fromDate, ToDate: toDate}

        result, err := mediatr.Send[*reportCommand.ExportStockMovementExcelRequest, *reportCommand.ExportStockMovementExcelResult](c.Context(), req)
        if err != nil {
            logger.Error("Failed to export stock movement Excel", "error", err)
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "Failed to export stock movements",
            })
        }

        c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
        c.Set("Content-Disposition", "attachment; filename="+result.Filename)
        return c.Send(result.Data)
    }
}
