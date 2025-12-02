package report

import (
	"log/slog"
	"mini-erp-backend/api/service/report/query"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

// PurchaseSummary
//
// 	@Summary		Get purchase summary
// 	@Description	Get purchase order summary by month including aggregate totals
// 	@Tags			Report
// 	@Accept			json
// 	@Produce		json
// 	@Param			month	query	string	true	"Month (MM-YYYY)"
// 	@Success		200	{object}	query.PurchaseSummaryResult
// 	@Failure		400	{object}	fiber.Map
// 	@Failure		500	{object}	fiber.Map
// 	@Router			/api/v1/reports/purchase-summary [get]
func PurchaseSummary(logger *slog.Logger) fiber.Handler {
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

        req := &query.PurchaseSummaryRequest{Year: monthDate.Year(), Month: int(monthDate.Month())}

        result, err := mediatr.Send[*query.PurchaseSummaryRequest, *query.PurchaseSummaryResult](c.Context(), req)
        if err != nil {
            logger.Error("Failed to get purchase summary", "error", err)
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "Failed to retrieve purchase summary",
            })
        }

        return c.Status(fiber.StatusOK).JSON(result)
    }
}
