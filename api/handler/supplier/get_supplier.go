package supplier

import (
	"log/slog"
	"mini-erp-backend/api/service/supplier/query"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

// Supplier
//
//	@Summary		Get a supplier by ID
//	@Description	Get supplier details by supplier ID
//	@Tags			Supplier
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Supplier ID (UUID)"
//	@Success		200	{object}	model.Supplier
//	@Failure	400	{object}	api.ErrorResponse
//	@Failure	404	{object}	api.ErrorResponse
//	@Failure	500	{object}	api.ErrorResponse
//	@Router			/suppliers/{id} [get]
func Supplier(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		supplierId, err := uuid.Parse(idParam)
		if err != nil {
			logger.Error("Invalid supplier ID", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid supplier ID"})
		}

		req := query.SupplierRequest{SupplierId: supplierId}

		result, err := mediatr.Send[*query.SupplierRequest, *query.SupplierResult](c.Context(), &req)
		if err != nil {
			logger.Error("Failed to get supplier", "error", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Supplier not found"})
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}
