package purchase_order_item

import (
	"log/slog"
	"mini-erp-backend/api/service/purchase_order_item/query"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

// PurchaseOrderItems
//
//	@Summary		List purchase order items
//	@Description	Get all items for a purchase order
//	@Tags			PurchaseOrderItem
//	@Produce		json
//	@Param			po_id	path	string	true	"Purchase Order ID"
//	@Success		200	{array}		model.PurchaseOrderItem
//	@Failure		400	{object}	api.ErrorResponse
//	@Failure		404	{object}	api.ErrorResponse
//	@Failure		500	{object}	api.ErrorResponse
//	@Router			/purchase-orders/{po_id}/items [get]
func PurchaseOrderItems(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		poIdStr := c.Params("po_id")
		poId, err := uuid.Parse(poIdStr)
		if err != nil {
			logger.Error("Invalid purchase order ID", "po_id", poIdStr, "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid purchase order ID",
			})
		}

		req := &query.PurchaseOrderItemsRequest{
			PurchaseOrderId: poId,
		}

		result, err := mediatr.Send[*query.PurchaseOrderItemsRequest, interface{}](c.Context(), req)
		if err != nil {
			logger.Error("Failed to get purchase order items", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}
