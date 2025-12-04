package stocktransaction_handler

import (
	"log/slog"
	"mini-erp-backend/api/service/stock_transaction/command"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

// StockIn is a function to handle stock in transactions
//
//	@Summary		Stock In
//	@Description	Handle stock in for a product
//	@Tags			StockTransaction
//	@Accept			json
//	@Produce		json
//	@Param			request	body		command.StockInRequest	true	"Stock In Request"
//	@Success		201		{object}	command.StockInResult
//	@Failure		400		{object}	api.ErrorResponse	"Bad Request: Invalid input or insufficient stock"
//	@Failure		404		{object}	api.ErrorResponse	"Not Found: Product does not exist"
//	@Failure		500		{object}	api.ErrorResponse	"Internal Server Error"s
//	@Router			/stock/in [post]
func StockIn(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := command.StockInRequest{}

		if err := c.BodyParser(&request); err != nil {
			logger.Error("Failed to parse stock in request", slog.String("error", err.Error()))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// ตรวจสอบ quantity ห้ามเป็นค่าติดลบหรือศูนย์
		if request.Quantity <= 0 {
			logger.Error("Invalid quantity for stock in", slog.Int64("quantity", request.Quantity))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Quantity must be greater than 0 for stock IN operation",
			})
		}

		if request.ReferenceId == nil {
			logger.Error("ReferenceId is required for stock in")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "ReferenceId is required for stock IN operation",
			})
		}

		response, err := mediatr.Send[command.StockInRequest, *command.StockInResult](c.Context(), request)

		if err != nil {
			logger.Error("Failed to process stock in", slog.String("error", err.Error()))

			if strings.Contains(err.Error(), "not found") {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			if strings.Contains(err.Error(), "quantity") {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to process stock in",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(response)
	}
}