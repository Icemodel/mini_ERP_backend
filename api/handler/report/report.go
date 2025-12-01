package report

import (
	"log/slog"
	reportCommand "mini-erp-backend/api/service/report/command"
	"mini-erp-backend/api/service/report/query"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

// GetStockSummary
//
//	@Summary		Get stock summary
//	@Description	Get current stock summary including aggregates and low stock products
//	@Tags			Report
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	query.GetStockSummaryResult
//	@Failure		500	{object}	fiber.Map
//	@Router			/api/v1/reports/stock-summary [get]
func GetStockSummary(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := &query.GetStockSummaryRequest{}

		result, err := mediatr.Send[*query.GetStockSummaryRequest, *query.GetStockSummaryResult](c.Context(), req)
		if err != nil {
			logger.Error("Failed to get stock summary", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve stock summary",
			})
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}

// GetStockMovements
//
//	@Summary		Get stock movements
//	@Description	Get stock movements within a date range
//	@Tags			Report
//	@Accept			json
//	@Produce		json
//	@Param			from	query	string	true	"From date (DD-MM-YYYY)"
//	@Param			to		query	string	true	"To date (DD-MM-YYYY)"
//	@Success		200	{object}	query.GetStockMovementsResult
//	@Failure		400	{object}	fiber.Map
//	@Failure		500	{object}	fiber.Map
//	@Router			/api/v1/reports/stock-movements [get]
func GetStockMovements(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fromStr := c.Query("from")
		toStr := c.Query("to")

		if fromStr == "" || toStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "from and to query parameters are required (format: DD-MM-YYYY)",
			})
		}

		// Parse DD-MM-YYYY format (layout: 02-01-2006)
		fromDate, err := time.Parse("02-01-2006", fromStr)
		if err != nil {
			logger.Error("Invalid from date", "from", fromStr, "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid from date format (expected: DD-MM-YYYY)",
			})
		}

		// Parse DD-MM-YYYY format (layout: 02-01-2006)
		toDate, err := time.Parse("02-01-2006", toStr)
		if err != nil {
			logger.Error("Invalid to date", "to", toStr, "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid to date format (expected: DD-MM-YYYY)",
			})
		}

		// Set time to end of day for toDate
		toDate = toDate.Add(24*time.Hour - time.Second)

		req := &query.GetStockMovementsRequest{
			FromDate: fromDate,
			ToDate:   toDate,
		}

		result, err := mediatr.Send[*query.GetStockMovementsRequest, *query.GetStockMovementsResult](c.Context(), req)
		if err != nil {
			logger.Error("Failed to get stock movements", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve stock movements",
			})
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}

// GetPurchaseSummary
//
//	@Summary		Get purchase summary
//	@Description	Get purchase order summary by month including aggregate totals
//	@Tags			Report
//	@Accept			json
//	@Produce		json
//	@Param			month	query	string	true	"Month (MM-YYYY)"
//	@Success		200	{object}	query.GetPurchaseSummaryResult
//	@Failure		400	{object}	fiber.Map
//	@Failure		500	{object}	fiber.Map
//	@Router			/api/v1/reports/purchase-summary [get]
func GetPurchaseSummary(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		monthStr := c.Query("month")

		if monthStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "month query parameter is required (format: MM-YYYY)",
			})
		}

		// Parse MM-YYYY format (layout: 01-2006)
		monthDate, err := time.Parse("01-2006", monthStr)
		if err != nil {
			logger.Error("Invalid month format", "month", monthStr, "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid month format (expected: MM-YYYY)",
			})
		}

		req := &query.GetPurchaseSummaryRequest{
			Year:  monthDate.Year(),
			Month: int(monthDate.Month()),
		}

		result, err := mediatr.Send[*query.GetPurchaseSummaryRequest, *query.GetPurchaseSummaryResult](c.Context(), req)
		if err != nil {
			logger.Error("Failed to get purchase summary", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve purchase summary",
			})
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}

// ExportStockSummaryCSV
//
//	@Summary		Export stock summary to CSV
//	@Description	Export current stock summary to CSV file
//	@Tags			Report
//	@Produce		text/csv
//	@Success		200	{file}	file
//	@Failure		500	{object}	fiber.Map
//	@Router			/api/v1/reports/stock-summary/export [get]
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

// ExportStockMovementExcel
//
//	@Summary		Export stock movements to Excel
//	@Description	Export stock movements within a date range to Excel file
//	@Tags			Report
//	@Produce		application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
//	@Param			from	query	string	true	"From date (DD-MM-YYYY)"
//	@Param			to		query	string	true	"To date (DD-MM-YYYY)"
//	@Success		200	{file}	file
//	@Failure		400	{object}	fiber.Map
//	@Failure		500	{object}	fiber.Map
//	@Router			/api/v1/reports/stock-movements/export [get]
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

		req := &reportCommand.ExportStockMovementExcelRequest{
			FromDate: fromDate,
			ToDate:   toDate,
		}

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

// ExportPurchaseReportExcel
//
//	@Summary		Export purchase report to Excel
//	@Description	Export purchase order summary by month to Excel file
//	@Tags			Report
//	@Produce		application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
//	@Param			month	query	string	true	"Month (MM-YYYY)"
//	@Success		200	{file}	file
//	@Failure		400	{object}	fiber.Map
//	@Failure		500	{object}	fiber.Map
//	@Router			/api/v1/reports/purchase-summary/export [get]
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

		req := &reportCommand.ExportPurchaseReportExcelRequest{
			Year:  monthDate.Year(),
			Month: int(monthDate.Month()),
		}

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
