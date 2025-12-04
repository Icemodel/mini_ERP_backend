package command

import (
	"context"
	"errors"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdatePurchaseOrder struct {
	logger      *slog.Logger
	db          *gorm.DB
	PORepo      repository.PurchaseOrder
	ProductRepo repository.Product
}

type UpdatePurchaseOrderRequest struct {
	PurchaseOrderId uuid.UUID `json:"-" `
	SupplierId      uuid.UUID `json:"supplier_id" validate:"required"`
}

func NewUpdatePurchaseOrder(
	logger *slog.Logger,
	db *gorm.DB,
	poRepo repository.PurchaseOrder,
	productRepo repository.Product,
) *UpdatePurchaseOrder {
	return &UpdatePurchaseOrder{
		logger:      logger,
		db:          db,
		PORepo:      poRepo,
		ProductRepo: productRepo,
	}
}

func (h *UpdatePurchaseOrder) Handle(ctx context.Context, req *UpdatePurchaseOrderRequest) (interface{}, error) {
	// Begin transaction
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Find PO
	po, err := h.PORepo.Search(tx, map[string]interface{}{
		"purchase_order_id": req.PurchaseOrderId,
	}, "")
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Only DRAFT can be updated
	if po.Status != model.Draft {
		tx.Rollback()
		h.logger.Error("Cannot update purchase order - invalid status", "status", po.Status)
		return nil, errors.New("can only update draft purchase orders")
	}

	// Update PO
	po.SupplierId = req.SupplierId
	
	if err := h.PORepo.Update(tx, po); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		h.logger.Error("Failed to commit transaction", "error", err)
		return nil, err
	}

	h.logger.Info("Purchase order updated successfully", "po_id", req.PurchaseOrderId)
	return po, nil
}
