package command

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeletePurchaseOrderItem struct {
	logger     *slog.Logger
	db         *gorm.DB
	POItemRepo repository.PurchaseOrderItem
	PORepo     repository.PurchaseOrder
}

type DeletePurchaseOrderItemRequest struct {
	PurchaseOrderItemId uuid.UUID
	PurchaseOrderId     uuid.UUID
}
type DeletePurchaseOrderItemResult struct {
	Deleted             bool      `json:"deleted"`
	PurchaseOrderItemId uuid.UUID `json:"purchase_order_item_id"`
}

func NewDeletePurchaseOrderItem(
	logger *slog.Logger,
	db *gorm.DB,
	poItemRepo repository.PurchaseOrderItem,
	poRepo repository.PurchaseOrder,
) *DeletePurchaseOrderItem {
	return &DeletePurchaseOrderItem{
		logger:     logger,
		db:         db,
		POItemRepo: poItemRepo,
		PORepo:     poRepo,
	}
}

func (h *DeletePurchaseOrderItem) Handle(ctx context.Context, req *DeletePurchaseOrderItemRequest) (*DeletePurchaseOrderItemResult, error) {
	// Get existing item first to get purchase_order_id
	item_id := map[string]interface{}{
		"purchase_order_item_id": req.PurchaseOrderItemId,
	}

	item, err := h.POItemRepo.Search(h.db, item_id, "")
	if err != nil {
		h.logger.Error("Failed to find purchase order item", "error", err)
		return nil, err
	}

	// Begin transaction
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Verify PO is DRAFT
	po, err := h.PORepo.Search(tx, map[string]interface{}{
		"purchase_order_id": req.PurchaseOrderId,
	}, "")
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if po.Status != model.Draft {
		tx.Rollback()
		h.logger.Error("Cannot delete items from non-draft purchase order", "status", po.Status)
		return nil, gorm.ErrInvalidData
	}

	// Delete item
	if err := h.POItemRepo.Delete(tx, req.PurchaseOrderItemId); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := h.PORepo.Update(tx, po); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		h.logger.Error("Failed to commit transaction", "error", err)
		return nil, err
	}

	h.logger.Info("Purchase order item deleted", "item_id", req.PurchaseOrderItemId)
	return &DeletePurchaseOrderItemResult{
		Deleted:             true,
		PurchaseOrderItemId: req.PurchaseOrderItemId,
	}, nil
}
