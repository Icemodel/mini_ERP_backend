package stocktransaction_handler

import (
	"log/slog"
	"mini-erp-backend/api/service/stock_transaction/command"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

// StockOut is a function to handle stock out transactions
//
//	@Summary		Stock Out
//	@Description	Handle stock out for a product
//	@Tags			StockTransaction
//	@Accept			json
//	@Produce		json
//	@Param			request	body		command.StockOutRequest	true	"Stock Out Request"
//	@Success		201		{object}	command.StockOutResult
//	@Failure		400		{object}	api.ErrorResponse	"Bad Request: Invalid input or insufficient stock"
//	@Failure		404		{object}	api.ErrorResponse	"Not Found: Product does not exist"
//	@Failure		500		{object}	api.ErrorResponse	"Internal Server Error"
//	@Router			/stocks/out [post]
func StockOut(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := command.StockOutRequest{}

		if err := c.BodyParser(&request); err != nil {
			logger.Error("Failed to parse stock out request", slog.String("error", err.Error()))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// ตรวจสอบ quantity ห้ามเป็นค่าติดลบหรือศูนย์
		if request.Quantity <= 0 {
			logger.Error("Invalid quantity for stock out", slog.Int64("quantity", request.Quantity))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Quantity must be greater than 0 for stock OUT operation",
			})
		}

		response, err := mediatr.Send[command.StockOutRequest, *command.StockOutResult](c.Context(), request)

		if err != nil {
			logger.Error("Failed to process stock out", slog.String("error", err.Error()))

			if strings.Contains(err.Error(), "not found") {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			if strings.Contains(err.Error(), "insufficient") {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			if strings.Contains(err.Error(), "quantity") {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to process stock out",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(response)
	}
}
