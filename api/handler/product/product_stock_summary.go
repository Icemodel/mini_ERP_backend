package product_handler

import (
	"log/slog"
	"mini-erp-backend/api/service/product/query"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

// ProductStockSummary is a function to get product stock summary by its ID
//
//	@Summary		Get Product Stock Summary
//	@Description	Get product stock summary by its ID
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Product ID"
//	@Success		200	{object}	query.ProductStockSummaryResult
//	@Router			/products/{id}/stock-summary [get]
func ProductStockSummary(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		productIdParam := c.Params("id")
		productId, err := uuid.Parse(productIdParam)

		if err != nil {
			logger.Error("Invalid product ID", slog.String("error", err.Error()))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid product ID",
			})
		}

		request := query.ProductStockSummaryRequest{
			ProductId: productId,
		}

		response, err := mediatr.Send[query.ProductStockSummaryRequest, *query.ProductStockSummaryResult](c.Context(), request)

		if err != nil {
			logger.Error("Failed to get product stock summary", slog.String("error", err.Error()))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get product stock summary",
			})
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}