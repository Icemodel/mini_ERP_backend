package purchase_order

import (
	"log/slog"
	"mini-erp-backend/api/service/purchase_order/command"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

// CreatePurchaseOrder
//
//	@Summary		Create a new purchase order
//	@Description	Create a new purchase order with items
//	@Tags			PurchaseOrder
//	@Accept			json
//	@Produce		json
//	@Param			purchaseOrder	body	command.CreatePurchaseOrderRequest	true	"Purchase Order information"
//	@Success		201	{object}	model.PurchaseOrder
//	@Failure		400	{object}	api.ErrorResponse
//	@Failure		500	{object}	api.ErrorResponse
//	@Router			/purchase-orders [post]
func CreatePurchaseOrder(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req command.CreatePurchaseOrderRequest

		err := c.BodyParser(&req)
		if err != nil {
			logger.Error("Failed to parse request body", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		result, err := mediatr.Send[*command.CreatePurchaseOrderRequest, interface{}](c.Context(), &req)
		if err != nil {
			logger.Error("Failed to create purchase order", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(result)
	}
}