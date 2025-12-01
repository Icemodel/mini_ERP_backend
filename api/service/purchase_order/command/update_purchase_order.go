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
	logger *slog.Logger
	db     *gorm.DB
	PORepo repository.PurchaseOrderRepository
}

type UpdatePurchaseOrderRequest struct {
	PurchaseOrderId uuid.UUID 
	SupplierId      uuid.UUID                 `json:"supplier_id" validate:"required"`
	Items           []UpdatePurchaseOrderItem `json:"items" validate:"required,min=1,dive"`
}

type UpdatePurchaseOrderItem struct {
	ProductId uuid.UUID `json:"product_id" validate:"required"`
	Quantity  uint64    `json:"quantity" validate:"required,min=1"`
	Price     float64   `json:"price" validate:"required,min=0"`
}


func NewUpdatePurchaseOrderHandler(
	logger *slog.Logger,
	db *gorm.DB,
	poRepo repository.PurchaseOrderRepository,
) *UpdatePurchaseOrder {
	return &UpdatePurchaseOrder{
		logger: logger,
		db:     db,
		PORepo: poRepo,
	}
}

func (h *UpdatePurchaseOrder) Handle(ctx context.Context, req *UpdatePurchaseOrderRequest) (interface{}, error) {
	// Find PO
	po, err := h.PORepo.FindById(h.db, req.PurchaseOrderId)
	if err != nil {
		return nil, err
	}

	// Only DRAFT can be updated
	if po.Status != model.Draft {
		h.logger.Error("Cannot update purchase order - invalid status", "status", po.Status)
		return nil, errors.New("can only update draft purchase orders")
	}

	// Begin transaction
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Calculate total amount from items
	var totalAmount uint64
	for _, it := range req.Items {
		totalAmount += uint64(it.Price * float64(it.Quantity))
	}

	// Update PO
	po.SupplierId = req.SupplierId
	po.TotalAmount = totalAmount

	if err := h.PORepo.Update(tx, po); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Delete existing items and replace with new items
	if err := h.PORepo.DeleteItemsByPOId(tx, po.PurchaseOrderId); err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, it := range req.Items {
		item := &model.PurchaseOrderItem{
			PurchaseOrderItemId: uuid.New(),
			PurchaseOrderId:     po.PurchaseOrderId,
			ProductId:           it.ProductId,
			Quantity:            it.Quantity,
			Price:               it.Price,
		}
		if err := h.PORepo.CreateItem(tx, item); err != nil {
			tx.Rollback()
			return nil, err
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

	h.logger.Info("Purchase order updated successfully", "po_id", req.PurchaseOrderId)
	return updatedPO, nil
}
