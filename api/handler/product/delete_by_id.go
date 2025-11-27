package product_handler

import (
	"log/slog"
	"mini-erp-backend/api/service/product/command"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

func DeleteById(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		productIdParam := c.Params("id")
		productId, err := uuid.Parse(productIdParam)

		if err != nil {
			logger.Error("Invalid product ID", slog.String("error", err.Error()))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid product ID",
			})
		}

		request := command.DeleteByIdRequest{
			ProductId: productId,
		}

		response, err := mediatr.Send[command.DeleteByIdRequest, *command.DeleteByIdResult](c.Context(), request)

		if err != nil {
			logger.Error("Failed to delete product", slog.String("error", err.Error()))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete product",
			})
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}
