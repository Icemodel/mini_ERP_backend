package purchase_order_item

import (
	"log/slog"
	"mini-erp-backend/api/service/purchase_order_item/command"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

// CreatePurchaseOrderItem
//
//	@Summary		Add item to purchase order
//	@Description	Add a new item to a draft purchase order
//	@Tags			PurchaseOrderItem
//	@Accept			json
//	@Produce		json
//	@Param			item	body	command.CreatePurchaseOrderItemRequest	true	"Item information"
//	@Success		201	{object}	model.PurchaseOrderItem
//	@Failure		400	{object}	api.ErrorResponse
//	@Failure		404	{object}	api.ErrorResponse
//	@Failure		500	{object}	api.ErrorResponse
//	@Router			/purchase-order-items [post]
func CreatePurchaseOrderItem(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req command.CreatePurchaseOrderItemRequest
		if err := c.BodyParser(&req); err != nil {
			logger.Error("Failed to parse request body", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		result, err := mediatr.Send[*command.CreatePurchaseOrderItemRequest, *command.CreatePurchaseOrderItemResult](c.Context(), &req)
		if err != nil {
			logger.Error("Failed to create purchase order item", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(result)
	}
}
