package product_handler

import (
	"log/slog"
	"mini-erp-backend/api/service/product/query"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

// ProductById is a function to get product by id
//
//	@Summary		Get Product by ID
//	@Description	Get product by ID
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	query.ProductResult
//	@Router			/products/{id} [get]
//
//	@param			id	path	string	true	"Product ID"
func ProductById(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		productIdParam := c.Params("id")
		productId, err := uuid.Parse(productIdParam)

		if err != nil {
			logger.Error("Invalid product ID", slog.String("error", err.Error()))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid product ID",
			})
		}

		request := query.ProductByIdRequest{
			ProductId: productId,
		}

		response, err := mediatr.Send[query.ProductByIdRequest, *query.ProductResult](c.Context(), request)
		if err != nil {
			logger.Error("Failed to get product by id", slog.String("error", err.Error()))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get product by id",
			})
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}