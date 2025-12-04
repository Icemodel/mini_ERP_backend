package purchase_order_item

import (
	"log/slog"
	"mini-erp-backend/api/service/purchase_order_item/query"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

// PurchaseOrderItem
//
//	@Summary		Get purchase order item
//	@Description	Get details of a specific purchase order item
//	@Tags			PurchaseOrderItem
//	@Produce		json
//	@Param			item_id	path	string	true	"Purchase Order Item ID"
//	@Success		200	{object}	model.PurchaseOrderItem
//	@Failure		400	{object}	api.ErrorResponse
//	@Failure		404	{object}	api.ErrorResponse
//	@Failure		500	{object}	api.ErrorResponse
//	@Router			/purchase-order-items/item/{item_id} [get]
func PurchaseOrderItem(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		itemIdStr := c.Params("item_id")
		itemId, err := uuid.Parse(itemIdStr)
		if err != nil {
			logger.Error("Invalid item ID", "item_id", itemIdStr, "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid item ID",
			})
		}

		req := &query.PurchaseOrderItemRequest{
			PurchaseOrderItemId: itemId,
		}

		result, err := mediatr.Send[*query.PurchaseOrderItemRequest, interface{}](c.Context(), req)
		if err != nil {
			logger.Error("Failed to get purchase order item", "error", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}
