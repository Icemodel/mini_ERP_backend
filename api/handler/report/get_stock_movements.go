package report

import (
	"log/slog"
	"mini-erp-backend/api/service/report/query"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

// StockMovements
//
//	@Summary		Get stock movements
//	@Description	Get stock movements within a date range
//	@Tags			Report
//	@Accept			json
//	@Produce		json
//	@Param			from	query		string	true	"From date (DD-MM-YYYY)"
//	@Param			to		query		string	true	"To date (DD-MM-YYYY)"
//	@Success		200		{object}	query.StockMovementsResult
//	@Failure		400		{object}	fiber.Map
//	@Failure		500		{object}	fiber.Map
//	@Router			/api/v1/reports/stock-movements [get]
func StockMovements(logger *slog.Logger) fiber.Handler {
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

		req := &query.StockMovementsRequest{
			FromDate: fromDate,
			ToDate:   toDate,
		}

		result, err := mediatr.Send[*query.StockMovementsRequest, *query.StockMovementsResult](c.Context(), req)
		if err != nil {
			logger.Error("Failed to get stock movements", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve stock movements",
			})
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}
