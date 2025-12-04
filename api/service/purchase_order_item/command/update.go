package command

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdatePurchaseOrderItem struct {
	logger     *slog.Logger
	db         *gorm.DB
	POItemRepo repository.PurchaseOrderItem
	PORepo     repository.PurchaseOrder
}

type UpdatePurchaseOrderItemRequest struct {
	PurchaseOrderItemId uuid.UUID `json:"-" swaggerignore:"true"`
	Quantity            uint64 `json:"quantity" validate:"required,min=1"`
}

func NewUpdatePurchaseOrderItem(
	logger *slog.Logger,
	db *gorm.DB,
	poItemRepo repository.PurchaseOrderItem,
	poRepo repository.PurchaseOrder,
) *UpdatePurchaseOrderItem {
	return &UpdatePurchaseOrderItem{
		logger:     logger,
		db:         db,
		POItemRepo: poItemRepo,
		PORepo:     poRepo,
	}
}

func (h *UpdatePurchaseOrderItem) Handle(ctx context.Context, req *UpdatePurchaseOrderItemRequest) (interface{}, error) {
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get existing item first to get purchase_order_id
	item_id := map[string]interface{}{
		"purchase_order_item_id": req.PurchaseOrderItemId,
	}
	item, err := h.POItemRepo.Search(tx, item_id, "")
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Verify PO is DRAFT
	po_id := map[string]interface{}{
		"purchase_order_id": item.PurchaseOrderId,
	}
	po, err := h.PORepo.Search(tx, po_id, "")
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if po.Status != model.Draft {
		tx.Rollback()
		h.logger.Error("Cannot update items in non-draft purchase order", "status", po.Status)
		return nil, gorm.ErrInvalidData
	}

	// Update quantity
	item.Quantity = req.Quantity

	if err := h.POItemRepo.Update(tx, req.PurchaseOrderItemId, item); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		h.logger.Error("Failed to commit transaction", "error", err)
		return nil, err
	}

	h.logger.Info("Purchase order item updated", "item_id", req.PurchaseOrderItemId)
	return item, nil
}