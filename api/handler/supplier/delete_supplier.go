package supplier

import (
	"log/slog"
	"mini-erp-backend/api/service/supplier/command"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

// DeleteSupplier
//
//	@Summary		Delete a supplier
//	@Description	Delete a supplier by ID
//	@Tags			Supplier
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Supplier ID (UUID)"
//	@Success		200
//	@Failure		400	{object}	api.ErrorResponse
//	@Failure		404	{object}	api.ErrorResponse
//	@Failure		500	{object}	api.ErrorResponse
//	@Router			/suppliers/{id} [delete]
func DeleteSupplier(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		supplierId, err := uuid.Parse(idParam)
		if err != nil {
			logger.Error("Invalid supplier ID", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid supplier ID"})
		}

		req := command.DeleteSupplierRequest{SupplierId: supplierId}

		result, err := mediatr.Send[*command.DeleteSupplierRequest, *command.DeleteSupplierResult](c.Context(), &req)
		if err != nil {
			logger.Error("Failed to delete supplier", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}
