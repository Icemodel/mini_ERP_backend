package stocktransaction_handler

import (
	"log/slog"
	"mini-erp-backend/api/service/stock_transaction/command"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

func StockAdjust(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := command.StockAdjustRequest{}

		if err := c.BodyParser(&request); err != nil {
			logger.Error("Failed to parse stock adjust request", slog.String("error", err.Error()))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// ตรวจสอบ quantity ห้ามเป็นศูนย์ (ADJUST สามารถเป็นลบได้เพื่อปรับลด)
		if request.Quantity == 0 {
			logger.Error("Invalid quantity for stock adjust", slog.Int64("quantity", request.Quantity))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Quantity cannot be 0 for stock ADJUST operation",
			})
		}

		// ตรวจสอบ reason ต้องไม่ว่าง (required สำหรับ ADJUST)
		if strings.TrimSpace(request.Reason) == "" {
			logger.Error("Reason is required for stock adjust")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Reason is required for stock ADJUST operation",
			})
		}

		response, err := mediatr.Send[command.StockAdjustRequest, *command.StockAdjustResult](c.Context(), request)

		if err != nil {
			logger.Error("Failed to process stock adjust", slog.String("error", err.Error()))

			if strings.Contains(err.Error(), "not found") {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			if strings.Contains(err.Error(), "reason") || strings.Contains(err.Error(), "quantity") || strings.Contains(err.Error(), "negative") {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to process stock adjust",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(response)
	}
}
