package purchase_order_item

import (
	"log/slog"
	"mini-erp-backend/api/service/purchase_order_item/command"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

// DeletePurchaseOrderItem
//
//	@Summary		Delete purchase order item
//	@Description	Remove an item from a draft purchase order
//	@Tags			PurchaseOrderItem
//	@Produce		json
//	@Param			item_id	path	string	true	"Purchase Order Item ID"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		400	{object}	api.ErrorResponse
//	@Failure		404	{object}	api.ErrorResponse
//	@Failure		500	{object}	api.ErrorResponse
//	@Router			/purchase-order-items/{item_id} [delete]
func DeletePurchaseOrderItem(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		itemIdStr := c.Params("item_id")
		itemId, err := uuid.Parse(itemIdStr)
		if err != nil {
			logger.Error("Invalid item ID", "item_id", itemIdStr, "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid item ID",
			})
		}

		req := &command.DeletePurchaseOrderItemRequest{
			PurchaseOrderItemId: itemId,
		}

		result, err := mediatr.Send[*command.DeletePurchaseOrderItemRequest, interface{}](c.Context(), req)
		if err != nil {
			logger.Error("Failed to delete purchase order item", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}