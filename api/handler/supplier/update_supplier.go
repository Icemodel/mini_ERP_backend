package supplier

import (
	"log/slog"
	"mini-erp-backend/api/service/supplier/command"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

// UpdateSupplier
//
//	@Summary		Update a supplier
//	@Description	Update supplier information by ID
//	@Tags			Supplier
//	@Accept			json
//	@Produce		json
//	@Param			id			path	string	true	"Supplier ID (UUID)"
//	@Param			supplier	body	command.UpdateSupplierRequest	true	"Updated supplier information"
//	@Success		200
//	@Failure		400	{object}	api.ErrorResponse
//	@Failure		404	{object}	api.ErrorResponse
//	@Failure		500	{object}	api.ErrorResponse
//	@Router			/suppliers/{id} [put]
func UpdateSupplier(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		supplierId, err := uuid.Parse(idParam)
		if err != nil {
			logger.Error("Invalid supplier ID", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid supplier ID"})
		}

		var req command.UpdateSupplierRequest
		err = c.BodyParser(&req)
		if err != nil {
			logger.Error("Failed to parse request body", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		req.SupplierId = supplierId

		_, err = mediatr.Send[*command.UpdateSupplierRequest, interface{}](c.Context(), &req)
		if err != nil {
			logger.Error("Failed to update supplier", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.SendStatus(fiber.StatusOK)
	}
}
