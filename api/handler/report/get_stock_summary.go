package report

import (
	"log/slog"
	"mini-erp-backend/api/service/report/query"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

// GetStockSummary
//
// 	@Summary		Get stock summary
// 	@Description	Get current stock summary including aggregates and low stock products
// 	@Tags			Report
// 	@Accept			json
// 	@Produce		json
// 	@Success		200	{object}	query.GetStockSummaryResult
// 	@Failure		500	{object}	fiber.Map
// 	@Router			/api/v1/reports/stock-summary [get]
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
