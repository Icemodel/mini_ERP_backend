package purchase_order

import (
	"log/slog"
	"mini-erp-backend/api/service/purchase_order/command"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

// UpdatePurchaseOrder
//
//	@Summary		Update a purchase order
//	@Description	Update purchase order information by ID
//	@Tags			PurchaseOrder
//	@Accept			json
//	@Produce		json
//	@Param			id				path		string								true	"Purchase Order ID (UUID)"
//	@Param			purchaseOrder	body		command.UpdatePurchaseOrderRequest	true	"Updated purchase order information"
//	@Success		200				{object}	model.PurchaseOrder
//	@Failure		400				{object}	fiber.Map
//	@Failure		500				{object}	fiber.Map
//	@Router			/api/v1/purchase-orders/{id} [put]
func UpdatePurchaseOrder(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		poIdStr := c.Params("id")
		poId, err := uuid.Parse(poIdStr)
		if err != nil {
			logger.Error("Invalid purchase order ID", "id", poIdStr, "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid purchase order ID",
			})
		}

		var req command.UpdatePurchaseOrderRequest
		err = c.BodyParser(&req)
		if err != nil {
			logger.Error("Failed to parse request body", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		req.PurchaseOrderId = poId

		result, err := mediatr.Send[*command.UpdatePurchaseOrderRequest, interface{}](c.Context(), &req)
		if err != nil {
			logger.Error("Failed to update purchase order", "po_id", poId, "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}
