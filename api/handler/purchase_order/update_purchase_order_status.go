package purchase_order

import (
	"log/slog"
	"mini-erp-backend/api/service/purchase_order/command"
	"mini-erp-backend/model"

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
//	@Param			id		status		path	string	true	"Purchase Order ID (UUID)"
//	@Param			status	body		object	true	"Status update"
//	@Success		200		{object}	model.PurchaseOrder
//	@Failure		400		{object}	fiber.Map
//	@Failure		500		{object}	fiber.Map
//	@Router			/api/v1/purchase-orders/{id}/status [put]
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

		var body struct {
			Status model.PurchaseOrderStatus `json:"status"`
		}
		err = c.BodyParser(&body)
		if err != nil {
			logger.Error("Failed to parse request body", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		req := &command.UpdatePOStatusRequest{
			PurchaseOrderId: poId,
			Status:          body.Status,
		}

		result, err := mediatr.Send[*command.UpdatePOStatusRequest, interface{}](c.Context(), req)
		if err != nil {
			logger.Error("Failed to update purchase order status", "po_id", poId, "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}
