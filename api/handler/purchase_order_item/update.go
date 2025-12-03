package purchase_order_item

import (
	"log/slog"
	"mini-erp-backend/api/service/purchase_order_item/command"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

// UpdatePurchaseOrderItem
//
//	@Summary		Update purchase order item
//	@Description	Update quantity of a purchase order item (draft only)
//	@Tags			PurchaseOrderItem
//	@Accept			json
//	@Produce		json
//	@Param			po_id	path	string	true	"Purchase Order ID"
//	@Param			item_id	path	string	true	"Purchase Order Item ID"
//	@Param			item	body	command.UpdatePurchaseOrderItemRequest	true	"Updated item information"
//	@Success		200	{object}	model.PurchaseOrderItem
//	@Failure		400	{object}	api.ErrorResponse
//	@Failure		404	{object}	api.ErrorResponse
//	@Failure		500	{object}	api.ErrorResponse
//	@Router			/purchase-orders/{po_id}/items/{item_id} [put]
func UpdatePurchaseOrderItem(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		poIdStr := c.Params("po_id")
		poId, err := uuid.Parse(poIdStr)
		if err != nil {
			logger.Error("Invalid purchase order ID", "po_id", poIdStr, "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid purchase order ID",
			})
		}

		itemIdStr := c.Params("item_id")
		itemId, err := uuid.Parse(itemIdStr)
		if err != nil {
			logger.Error("Invalid item ID", "item_id", itemIdStr, "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid item ID",
			})
		}

		var req command.UpdatePurchaseOrderItemRequest
		if err := c.BodyParser(&req); err != nil {
			logger.Error("Failed to parse request body", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		req.PurchaseOrderId = poId
		req.PurchaseOrderItemId = itemId

		result, err := mediatr.Send[*command.UpdatePurchaseOrderItemRequest, interface{}](c.Context(), &req)
		if err != nil {
			logger.Error("Failed to update purchase order item", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}
