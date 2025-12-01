package purchase_order

import (
	"log/slog"
	"mini-erp-backend/api/service/purchase_order/command"
	"mini-erp-backend/api/service/purchase_order/query"
	"mini-erp-backend/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

// GetAllPurchaseOrders
//
//	@Summary		Get all purchase orders
//	@Description	Retrieve a list of all purchase orders
//	@Tags			PurchaseOrder
//	@Accept			json
//	@Produce		json
//	@Param			status		query	string	false	"Filter by status"
//	@Param			order_by	query	string	false	"Order by field"
//	@Success		200	{array}	model.PurchaseOrder
//	@Failure		500	{object}	fiber.Map
//	@Router			/api/v1/purchase-orders [get]
func GetAllPurchaseOrders(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req query.GetAllPurchaseOrdersRequest
		
		if statusStr := c.Query("status"); statusStr != "" {
			status := model.PurchaseOrderStatus(statusStr)
			req.Status = &status
		}
		req.OrderBy = c.Query("order_by")

		result, err := mediatr.Send[*query.GetAllPurchaseOrdersRequest, *query.GetAllPurchaseOrdersResult](c.Context(), &req)
		if err != nil {
			logger.Error("Failed to get all purchase orders", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve purchase orders",
			})
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}

// GetPurchaseOrder
//
//	@Summary		Get a purchase order by ID
//	@Description	Get purchase order details by ID
//	@Tags			PurchaseOrder
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Purchase Order ID (UUID)"
//	@Success		200	{object}	model.PurchaseOrder
//	@Failure		400	{object}	fiber.Map
//	@Failure		404	{object}	fiber.Map
//	@Failure		500	{object}	fiber.Map
//	@Router			/api/v1/purchase-orders/{id} [get]
func GetPurchaseOrder(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		poIdStr := c.Params("id")
		poId, err := uuid.Parse(poIdStr)
		if err != nil {
			logger.Error("Invalid purchase order ID", "id", poIdStr, "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid purchase order ID",
			})
		}

		req := &query.GetPurchaseOrderRequest{
			PurchaseOrderId: poId,
		}

		result, err := mediatr.Send[*query.GetPurchaseOrderRequest, *query.GetPurchaseOrderResult](c.Context(), req)
		if err != nil {
			logger.Error("Failed to get purchase order", "po_id", poId, "error", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Purchase order not found",
			})
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}

// CreatePurchaseOrder
//
//	@Summary		Create a new purchase order
//	@Description	Create a new purchase order with items
//	@Tags			PurchaseOrder
//	@Accept			json
//	@Produce		json
//	@Param			purchaseOrder	body	command.CreatePurchaseOrderRequest	true	"Purchase Order information"
//	@Success		201	{object}	model.PurchaseOrder
//	@Failure		400	{object}	fiber.Map
//	@Failure		500	{object}	fiber.Map
//	@Router			/api/v1/purchase-orders [post]
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

// UpdatePurchaseOrder
//
//	@Summary		Update a purchase order
//	@Description	Update purchase order information by ID
//	@Tags			PurchaseOrder
//	@Accept			json
//	@Produce		json
//	@Param			id				path	string								true	"Purchase Order ID (UUID)"
//	@Param			purchaseOrder	body	command.UpdatePurchaseOrderRequest	true	"Updated purchase order information"
//	@Success		200	{object}	model.PurchaseOrder
//	@Failure		400	{object}	fiber.Map
//	@Failure		500	{object}	fiber.Map
//	@Router			/api/v1/purchase-orders/{id} [put]
func UpdatePurchaseOrder(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		poIdStr := c.Params("id")
		poId, err := uuid.Parse(poIdStr)
		if err != nil {
			logger.Error("Invalid purchase order ID", "id", poIdStr, "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid purchase order ID",
			})
		}

		var req command.UpdatePurchaseOrderRequest
		err = c.BodyParser(&req)
		if err != nil {
			logger.Error("Failed to parse request body", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		req.PurchaseOrderId = poId

		result, err := mediatr.Send[*command.UpdatePurchaseOrderRequest, interface{}](c.Context(), &req)
		if err != nil {
			logger.Error("Failed to update purchase order", "po_id", poId, "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}

// UpdatePurchaseOrderStatus
//
//	@Summary		Update purchase order status
//	@Description	Update the status of a purchase order
//	@Tags			PurchaseOrder
//	@Accept			json
//	@Produce		json
//	@Param			id		status	path	string	true	"Purchase Order ID (UUID)"
//	@Param			status	body	object	true	"Status update"
//	@Success		200	{object}	model.PurchaseOrder
//	@Failure		400	{object}	fiber.Map
//	@Failure		500	{object}	fiber.Map
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
