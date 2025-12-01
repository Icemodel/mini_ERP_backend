package product_handler

import (
	"log/slog"
	"mini-erp-backend/api/service/product/command"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

func Update(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		productIdParam := c.Params("id")
		productId, err := uuid.Parse(productIdParam)

		if err != nil {
			logger.Error("Invalid product ID", slog.String("error", err.Error()))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid product ID",
			})
		}

		request := command.UpdateRequest{
			ProductId: productId,
		}

		if err := c.BodyParser(&request); err != nil {
			logger.Error("Failed to parse update product request", slog.String("error", err.Error()))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		response, err := mediatr.Send[command.UpdateRequest, *command.UpdateResult](c.Context(), request)

		if err != nil {
			logger.Error("Failed to update product", slog.String("error", err.Error()))

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
				"error": "Failed to update product",
			})
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}
