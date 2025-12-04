package command

import (
	"context"
	"errors"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdatePOStatus struct {
	logger    *slog.Logger
	db        *gorm.DB
	PORepo    repository.PurchaseOrder
	POItemRepo repository.PurchaseOrderItem
	StockRepo repository.StockTransaction
}

type UpdatePOStatusRequest struct {
	PurchaseOrderId uuid.UUID                 `json:"-" swaggerignore:"true"`
	Status          model.PurchaseOrderStatus `json:"status" validate:"required"`
}

func NewUpdatePOStatus(
	logger *slog.Logger,
	db *gorm.DB,
	poRepo repository.PurchaseOrder,
	poItemRepo repository.PurchaseOrderItem,
	stockRepo repository.StockTransaction,
) *UpdatePOStatus {
	return &UpdatePOStatus{
		logger:    logger,
		db:        db,
		PORepo:    poRepo,
		POItemRepo: poItemRepo,
		StockRepo: stockRepo,
	}
}

func (h *UpdatePOStatus) Handle(ctx context.Context, req *UpdatePOStatusRequest) (interface{}, error) {
	// Begin transaction
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// If changing to RECEIVED, validate that items exist
	var items []*model.PurchaseOrderItem
	if req.Status == model.Received {
		var err error
		items, err = h.POItemRepo.Searches(tx, map[string]interface{}{"purchase_order_id": req.PurchaseOrderId}, "")
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if len(items) == 0 {
			tx.Rollback()
			h.logger.Error("Cannot receive purchase order without items", "po_id", req.PurchaseOrderId)
			return nil, errors.New("cannot receive purchase order without items")
		}
	}

	// Get PurchaseOrder to access CreatedBy
	po, err := h.PORepo.Search(tx, map[string]interface{}{"purchase_order_id": req.PurchaseOrderId}, "")
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Update status
	if err := h.PORepo.UpdateStatus(tx, req.PurchaseOrderId, req.Status); err != nil {
		tx.Rollback()
		return nil, err
	}

	// If status is RECEIVED, create Stock IN transactions
	if req.Status == model.Received {
		// Create Stock IN transactions for each item
		for _, item := range items {
			// Fetch latest stock transaction for the product
			latestStockTx, err := h.StockRepo.GetLatestByProduct(tx, item.ProductId)
			if err != nil {
				tx.Rollback()
				return nil, err
			}

			// Calculate new quantity 
			newQuantity := int64(item.Quantity)
			if latestStockTx != nil {
				newQuantity = latestStockTx.Quantity + int64(item.Quantity)
			}

			received := "Purchase Order Received"

			// Create stock transaction
			stockTx := &model.StockTransaction{
				StockTransactionId: uuid.New(),
				ProductId:          item.ProductId,
				Quantity:           newQuantity,
				Type:               "IN",
				Reason:             &received,
				ReferenceId:        &req.PurchaseOrderId,
				CreatedAt:          time.Now(),
				CreatedBy:          po.CreatedBy,
			}

			if err := h.StockRepo.Create(tx, stockTx); err != nil {
				tx.Rollback()
				return nil, err
			}

			if latestStockTx != nil {
				h.logger.Info("Stock transaction created with cumulative quantity",
					"product_id", item.ProductId,
					"po_item_quantity", item.Quantity,
					"previous_quantity", latestStockTx.Quantity,
					"new_quantity", newQuantity,
					"po_id", req.PurchaseOrderId)
			} else {
				h.logger.Info("Stock transaction created (first)",
					"product_id", item.ProductId,
					"quantity", newQuantity,
					"po_id", req.PurchaseOrderId)
			}
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		h.logger.Error("Failed to commit transaction", "error", err)
		return nil, err
	}

	h.logger.Info("Purchase order status updated", "po_id", req.PurchaseOrderId, "status", req.Status)
	return map[string]interface{}{
		"purchase_order_id": req.PurchaseOrderId,
		"status":            req.Status,
	}, nil
}