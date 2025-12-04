package purchase_order_item

import (
	"log/slog"
	"mini-erp-backend/api/service/purchase_order_item/query"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

// AllPurchaseOrderItems
//
//	@Summary		List all purchase order items
//	@Description	Get all purchase order items across all purchase orders
//	@Tags			PurchaseOrderItem
//	@Produce		json
//	@Success		200	{array}		model.PurchaseOrderItem
//	@Failure		500	{object}	api.ErrorResponse
//	@Router			/purchase-order-items [get]
func AllPurchaseOrderItems(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := &query.AllPurchaseOrderItemsRequest{}

		result, err := mediatr.Send[*query.AllPurchaseOrderItemsRequest, *query.AllPurchaseOrderItemsResult](c.Context(), req)
		if err != nil {
			logger.Error("Failed to get all purchase order items", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}