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
	StockRepo repository.StockTransaction
}

type UpdatePOStatusRequest struct {
	PurchaseOrderId uuid.UUID               
	Status          model.PurchaseOrderStatus `json:"status" validate:"required"`
}

func NewUpdatePOStatusHandler(
	logger *slog.Logger,
	db *gorm.DB,
	poRepo repository.PurchaseOrder,
	stockRepo repository.StockTransaction,
) *UpdatePOStatus {
	return &UpdatePOStatus{
		logger:    logger,
		db:        db,
		PORepo:    poRepo,
		StockRepo: stockRepo,
	}
}

func (h *UpdatePOStatus) Handle(ctx context.Context, req *UpdatePOStatusRequest) (interface{}, error) {
	// Find PO
	po, err := h.PORepo.Search(h.db, map[string]interface{}{
		"purchase_order_id": req.PurchaseOrderId,
	}, "")
	if err != nil {
		return nil, err
	}

	// If changing to RECEIVED, validate that items exist
	if req.Status == model.Received {
		items, err := h.PORepo.SearchItemsByPurchaseOrderId(h.db, req.PurchaseOrderId)
		if err != nil {
			return nil, err
		}
		if len(items) == 0 {
			h.logger.Error("Cannot receive purchase order without items", "po_id", req.PurchaseOrderId)
			return nil, errors.New("cannot receive purchase order without items")
		}
	}

	// Begin transaction
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update status
	if err := h.PORepo.UpdateStatus(tx, req.PurchaseOrderId, req.Status); err != nil {
		tx.Rollback()
		return nil, err
	}

	// If status is RECEIVED, create Stock IN transactions
	if req.Status == model.Received {
		// Get PO items
		items, err := h.PORepo.SearchItemsByPurchaseOrderId(h.db, req.PurchaseOrderId)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Create Stock IN transactions for each item
		for _, item := range items {
			stockTx := &model.StockTransaction{
				StockTransactionId: uuid.New(),
				ProductId:          item.ProductId,
				Quantity:           int64(item.Quantity), // IN = positive
				Type:               "IN",
				Reason:             stringPtr("Purchase Order Received"),
				ReferenceId:        &req.PurchaseOrderId,
				CreatedAt:          time.Now(),
				CreatedBy:          po.CreatedBy.String(),
			}

			if err := h.StockRepo.Create(tx, stockTx); err != nil {
				tx.Rollback()
				return nil, err
			}

			h.logger.Info("Stock transaction created",
				"product_id", item.ProductId,
				"quantity", item.Quantity,
				"po_id", req.PurchaseOrderId)
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		h.logger.Error("Failed to commit transaction", "error", err)
		return nil, err
	}

	// Fetch updated PO
	updatedPO, err := h.PORepo.Search(h.db, map[string]interface{}{
		"purchase_order_id": req.PurchaseOrderId,
	}, "")
	if err != nil {
		return nil, err
	}

	h.logger.Info("Purchase order status updated", "po_id", req.PurchaseOrderId, "status", req.Status)
	return updatedPO, nil
}

func stringPtr(s string) *string {
	return &s
}
