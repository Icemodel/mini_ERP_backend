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
	ProductId           uuid.UUID `json:"product_id" validate:"required"`
	Quantity            uint64    `json:"quantity" validate:"required,min=1"`
	Price               float64   `json:"price" validate:"required,min=0"`
}

type UpdatePurchaseOrderItemResult struct {
	PurchaseOrderItem model.PurchaseOrderItem `json:"purchase_order_item"`
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

func (h *UpdatePurchaseOrderItem) Handle(ctx context.Context, req *UpdatePurchaseOrderItemRequest) (*UpdatePurchaseOrderItemResult, error) {
	// Get existing item first to get purchase_order_id
	item_id := map[string]interface{}{
		"purchase_order_item_id": req.PurchaseOrderItemId,
	}
	item, err := h.POItemRepo.Search(h.db, item_id, "")
	if err != nil {
		h.logger.Error("Failed to find purchase order item", "error", err)
		return nil, err
	}

	if req.ProductId == uuid.Nil {
		h.logger.Error("Product ID is required")
		return nil, gorm.ErrInvalidData
	}

	if req.Quantity <= 0 {
		h.logger.Error("Quantity must be greater than zero")
		return nil, gorm.ErrInvalidData
	}

	if req.Price < 0 {
		h.logger.Error("Price cannot be negative")
		return nil, gorm.ErrInvalidData
	}

	// Begin transaction
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

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

	// Update fields
	item.ProductId = req.ProductId
	item.Quantity = req.Quantity
	item.Price = req.Price

	if err := h.POItemRepo.Update(tx, req.PurchaseOrderItemId, item); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		h.logger.Error("Failed to commit transaction", "error", err)
		return nil, err
	}
	
	return &UpdatePurchaseOrderItemResult{
		PurchaseOrderItem: *item,
	}, nil
}