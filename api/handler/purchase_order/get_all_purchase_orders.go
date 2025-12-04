package purchase_order

import (
	"log/slog"
	"mini-erp-backend/api/service/purchase_order/query"
	"mini-erp-backend/model"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

// AllPurchaseOrders
//
//	@Summary		Get all purchase orders
//	@Description	Retrieve a list of all purchase orders
//	@Tags			PurchaseOrder
//	@Accept			json
//	@Produce		json
//	@Param			status		query	string	false	"Filter by status"
//	@Param			order_by	query	string	false	"Order by field"
//	@Success		200	{array}	model.PurchaseOrder
//	@Failure		500	{object}	api.ErrorResponse
//	@Router			/purchase-orders [get]
func AllPurchaseOrders(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req query.AllPurchaseOrdersRequest

		if statusStr := c.Query("status"); statusStr != "" {
			status := model.PurchaseOrderStatus(statusStr)
			req.Status = &status
		}
		req.OrderBy = c.Query("order_by")

		result, err := mediatr.Send[*query.AllPurchaseOrdersRequest, *query.AllPurchaseOrdersResult](c.Context(), &req)
		if err != nil {
			logger.Error("Failed to get all purchase orders", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve purchase orders",
			})
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}