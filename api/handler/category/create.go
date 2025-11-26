package category_handler

import (
	"log/slog"
	"mini-erp-backend/api/service/category/command"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

func Create(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := command.CreateRequest{}

		if err := c.BodyParser(&request); err != nil {
			logger.Error("Failed to parse create category request", slog.String("error", err.Error()))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if request.Name == "" {
			logger.Error("Category name is required")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Category name is required",
			})
		}

		response, err := mediatr.Send[command.CreateRequest, *command.CreateResult](c.Context(), request)

		if err != nil {
			logger.Error("Failed to create category", slog.String("error", err.Error()))

			// ตรวจสอบว่าเป็น error ชื่อซ้ำหรือไม่
			if strings.Contains(err.Error(), "already exists") {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create category",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(response)
	}
}
