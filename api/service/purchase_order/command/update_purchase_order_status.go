package command

import (
	"context"
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
	PORepo    repository.PurchaseOrderRepository
	StockRepo repository.StockTransactionRepository
}

type UpdatePOStatusRequest struct {
	PurchaseOrderId uuid.UUID                 `json:"purchase_order_id" validate:"required"`
	Status          model.PurchaseOrderStatus `json:"status" validate:"required"`
	UpdatedBy       string                    `json:"updated_by"`
}

func NewUpdatePOStatusHandler(
	logger *slog.Logger,
	db *gorm.DB,
	poRepo repository.PurchaseOrderRepository,
	stockRepo repository.StockTransactionRepository,
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
	_, err := h.PORepo.FindById(h.db, req.PurchaseOrderId)
	if err != nil {
		return nil, err
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
		items, err := h.PORepo.FindItemsByPOId(h.db, req.PurchaseOrderId)
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
				CreatedBy:          req.UpdatedBy,
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
	updatedPO, err := h.PORepo.FindById(h.db, req.PurchaseOrderId)
	if err != nil {
		return nil, err
	}

	h.logger.Info("Purchase order status updated", "po_id", req.PurchaseOrderId, "status", req.Status)
	return updatedPO, nil
}

func stringPtr(s string) *string {
	return &s
}
