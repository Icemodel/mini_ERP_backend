package report

import (
	"log/slog"
	"mini-erp-backend/api/service/report/query"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

// StockSummary
//
//	@Summary		Get stock summary
//	@Description	Get current stock summary including aggregates and low stock products
//	@Tags			Report
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	query.StockSummaryResult
//	@Failure		500	{object}	api.ErrorResponse
//	@Router			/reports/stock-summary [get]
func StockSummary(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := &query.StockSummaryRequest{}

		result, err := mediatr.Send[*query.StockSummaryRequest, *query.StockSummaryResult](c.Context(), req)
		if err != nil {
			logger.Error("Failed to get stock summary", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve stock summary",
			})
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}