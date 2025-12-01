package category_handler

import (
	"log/slog"
	"mini-erp-backend/api/service/category/command"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

func DeleteById(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		categoryIdParam := c.Params("id")
		categoryId, err := uuid.Parse(categoryIdParam)
		if err != nil {
			logger.Error("Invalid category ID", slog.String("error", err.Error()))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid category ID",
			})
		}
		request := command.DeleteByIdRequest{
			CategoryId: categoryId,
		}

		response, err := mediatr.Send[command.DeleteByIdRequest, *command.DeleteByIdResult](c.Context(), request)

		if err != nil {
			logger.Error("Failed to delete category", slog.String("error", err.Error()))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete category",
			})
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}
