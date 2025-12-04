package purchase_order

import (
	"log/slog"
	"mini-erp-backend/api/service/purchase_order/command"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

// UpdatePurchaseOrderStatus
//
//	@Summary		Update purchase order status
//	@Description	Update the status of a purchase order
//	@Tags			PurchaseOrder
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string							true	"Purchase Order ID (UUID)"
//	@Param			request	body	command.UpdatePOStatusRequest	true	"Status update request"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{object}	api.ErrorResponse
//	@Failure		500		{object}	api.ErrorResponse
//	@Router			/purchase-orders/{id}/status [put]
func UpdatePurchaseOrderStatus(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		poIdStr := c.Params("id")
		poId, err := uuid.Parse(poIdStr)

		if err != nil {
			logger.Error("Invalid purchase order ID", "id", poIdStr, "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid purchase order ID",
			})
		}

		var req command.UpdatePOStatusRequest
		err = c.BodyParser(&req)
		if err != nil {
			logger.Error("Failed to parse request body", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		req.PurchaseOrderId = poId

		result, err := mediatr.Send[*command.UpdatePOStatusRequest, *command.UpdatePOStatusResult](c.Context(), &req)
		if err != nil {
			logger.Error("Failed to update purchase order status", "po_id", poId, "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}
