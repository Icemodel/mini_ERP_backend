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
	PurchaseOrderId uuid.UUID 
	SupplierId      uuid.UUID                 `json:"supplier_id" validate:"required"`
	Items           []UpdatePurchaseOrderItem `json:"items" validate:"required,min=1,dive"`
}

type UpdatePurchaseOrderItem struct {
	ProductId uuid.UUID `json:"product_id" validate:"required"`
	Quantity  uint64    `json:"quantity" validate:"required,min=1"`
}


func NewUpdatePurchaseOrderHandler(
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
	// Find PO
	po, err := h.PORepo.Search(h.db, map[string]interface{}{
		"purchase_order_id": req.PurchaseOrderId,
	}, "")
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

	// Fetch product prices and calculate total amount
	var totalAmount uint64
	itemPrices := make(map[uuid.UUID]float64)

	for _, it := range req.Items {
		// Fetch product to get cost price
		product, err := h.ProductRepo.Search(h.db, map[string]interface{}{
			"product_id": it.ProductId,
		}, "")
		if err != nil {
			tx.Rollback()
			h.logger.Error("Product not found", "product_id", it.ProductId, "error", err)
			return nil, err
		}

		// Store price snapshot
		itemPrices[it.ProductId] = product.CostPrice
		totalAmount += uint64(product.CostPrice * float64(it.Quantity))
	}

	// Update PO
	po.SupplierId = req.SupplierId
	po.TotalAmount = totalAmount

	if err := h.PORepo.Update(tx, po); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Delete existing items and replace with new items
	if err := h.PORepo.DeleteItemsByPurchaseOrderId(tx, po.PurchaseOrderId); err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, it := range req.Items {
		item := &model.PurchaseOrderItem{
			PurchaseOrderItemId: uuid.New(),
			PurchaseOrderId:     po.PurchaseOrderId,
			ProductId:           it.ProductId,
			Quantity:            it.Quantity,
			Price:               itemPrices[it.ProductId],
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
	updatedPO, err := h.PORepo.Search(h.db, map[string]interface{}{
		"purchase_order_id": req.PurchaseOrderId,
	}, "")
	if err != nil {
		return nil, err
	}

	h.logger.Info("Purchase order updated successfully", "po_id", req.PurchaseOrderId)
	return updatedPO, nil
}
