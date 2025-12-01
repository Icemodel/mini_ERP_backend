package category_handler

import (
	"log/slog"
	"mini-erp-backend/api/service/category/query"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

func CategoryById(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		categoryId, err := uuid.Parse(id)

		if err != nil {
			logger.Error("Invalid category ID", slog.String("error", err.Error()))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid category ID",
			})
		}

		request := query.CategoryByIdRequest{
			CategoryId: categoryId,
		}

		response, err := mediatr.Send[query.CategoryByIdRequest, *query.CategoryByIdResult](c.Context(), request)

		if err != nil {
			logger.Error("Failed to get category by id", slog.String("error", err.Error()))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get category by id",
			})
		}
		return c.Status(fiber.StatusOK).JSON(response)
	}
}
