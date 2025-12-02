package supplier

import (
	"log/slog"
	"mini-erp-backend/api/service/supplier/command"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

// CreateSupplier
//
//	@Summary		Create a new supplier
//	@Description	Create a new supplier with the provided information
//	@Tags			Supplier
//	@Accept			json
//	@Produce		json
//	@Param			supplier	body		command.CreateSupplierRequest	true	"Supplier information"
//	@Success		201			{object}	model.Supplier
//	@Failure		400			{object}	fiber.Map
//	@Failure		500			{object}	fiber.Map
//	@Router			/api/v1/suppliers [post]
func CreateSupplier(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req command.CreateSupplierRequest

		err := c.BodyParser(&req)
		if err != nil {
			logger.Error("Failed to parse request body", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		result, err := mediatr.Send[*command.CreateSupplierRequest, interface{}](c.Context(), &req)
		if err != nil {
			logger.Error("Failed to create supplier", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(result)
	}
}
