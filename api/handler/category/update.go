package category_handler

import (
	"log/slog"
	"mini-erp-backend/api/service/category/command"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

func Update(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := command.UpdateRequest{}

		if err := c.BodyParser(&request); err != nil {
			logger.Error("Failed to parse update category request", slog.String("error", err.Error()))
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

		response, err := mediatr.Send[command.UpdateRequest, *command.UpdateResult](c.Context(), request)

		if err != nil {
			logger.Error("Failed to update category", slog.String("error", err.Error()))

			// ตรวจสอบว่าเป็น error ชื่อซ้ำหรือไม่
			if strings.Contains(err.Error(), "already exists") {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			// ตรวจสอบว่าเป็น error ไม่พบข้อมูลหรือไม่
			if strings.Contains(err.Error(), "not found") {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update category",
			})
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}
