package report

import (
	"log/slog"
	"mini-erp-backend/api/service/report/query"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

// GetStockSummary
//
//	@Summary		Get stock summary
//	@Description	Get current stock summary for all products
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
//	@Param			from	query	string	true	"From date (YYYY-MM-DD)"
//	@Param			to		query	string	true	"To date (YYYY-MM-DD)"
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
				"error": "from and to query parameters are required (format: YYYY-MM-DD)",
			})
		}

		fromDate, err := time.Parse("2006-01-02", fromStr)
		if err != nil {
			logger.Error("Invalid from date", "from", fromStr, "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid from date format (expected: YYYY-MM-DD)",
			})
		}

		toDate, err := time.Parse("2006-01-02", toStr)
		if err != nil {
			logger.Error("Invalid to date", "to", toStr, "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid to date format (expected: YYYY-MM-DD)",
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
//	@Description	Get purchase order summary by month
//	@Tags			Report
//	@Accept			json
//	@Produce		json
//	@Param			month	query	string	true	"Month (YYYY-MM)"
//	@Success		200	{object}	query.GetPurchaseSummaryResult
//	@Failure		400	{object}	fiber.Map
//	@Failure		500	{object}	fiber.Map
//	@Router			/api/v1/reports/purchase-summary [get]
func GetPurchaseSummary(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		monthStr := c.Query("month")

		if monthStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "month query parameter is required (format: YYYY-MM)",
			})
		}

		// Parse YYYY-MM format
		monthDate, err := time.Parse("2006-01", monthStr)
		if err != nil {
			logger.Error("Invalid month format", "month", monthStr, "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid month format (expected: YYYY-MM)",
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
