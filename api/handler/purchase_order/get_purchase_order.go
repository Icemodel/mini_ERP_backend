package purchase_order

import (
	"log/slog"
	"mini-erp-backend/api/service/purchase_order/query"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

// PurchaseOrder
//
// 	@Summary		Get a purchase order by ID
// 	@Description	Get purchase order details by ID
// 	@Tags			PurchaseOrder
// 	@Accept			json
// 	@Produce		json
// 	@Param			id	path	string	true	"Purchase Order ID (UUID)"
// 	@Success		200	{object}	model.PurchaseOrder
// 	@Failure		400	{object}	api.ErrorResponse
// 	@Failure		404	{object}	api.ErrorResponse
// 	@Failure		500	{object}	api.ErrorResponse
// 	@Router			/api/v1/purchase-orders/{id} [get]
func PurchaseOrder(logger *slog.Logger) fiber.Handler {
    return func(c *fiber.Ctx) error {
        poIdStr := c.Params("id")
        poId, err := uuid.Parse(poIdStr)
        if err != nil {
            logger.Error("Invalid purchase order ID", "id", poIdStr, "error", err)
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": "Invalid purchase order ID",
            })
        }

        req := &query.PurchaseOrderRequest{
            PurchaseOrderId: poId,
        }

        result, err := mediatr.Send[*query.PurchaseOrderRequest, *query.PurchaseOrderResult](c.Context(), req)
        if err != nil {
            logger.Error("Failed to get purchase order", "po_id", poId, "error", err)
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "error": "Purchase order not found",
            })
        }

        return c.Status(fiber.StatusOK).JSON(result)
    }
}
