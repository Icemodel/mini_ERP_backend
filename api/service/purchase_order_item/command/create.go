package command

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreatePurchaseOrderItem struct {
	logger      *slog.Logger
	db          *gorm.DB
	POItemRepo  repository.PurchaseOrderItem
	PORepo      repository.PurchaseOrder
	ProductRepo repository.Product
}

type CreatePurchaseOrderItemRequest struct {
	PurchaseOrderId uuid.UUID `json:"purchase_order_id" validate:"required"`
	ProductId       uuid.UUID `json:"product_id" validate:"required"`
	Quantity        uint64    `json:"quantity" validate:"required,min=1"`
}

func NewCreatePurchaseOrderItem(
	logger *slog.Logger,
	db *gorm.DB,
	poItemRepo repository.PurchaseOrderItem,
	poRepo repository.PurchaseOrder,
	productRepo repository.Product,
) *CreatePurchaseOrderItem {
	return &CreatePurchaseOrderItem{
		logger:      logger,
		db:          db,
		POItemRepo:  poItemRepo,
		PORepo:      poRepo,
		ProductRepo: productRepo,
	}
}

func (h *CreatePurchaseOrderItem) Handle(ctx context.Context, req *CreatePurchaseOrderItemRequest) (interface{}, error) {
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Verify PO exists and is DRAFT
	po, err := h.PORepo.Search(tx, map[string]interface{}{
		"purchase_order_id": req.PurchaseOrderId,
	}, "")
	if err != nil {
		tx.Rollback()
		h.logger.Error("Purchase order not found", "po_id", req.PurchaseOrderId, "error", err)
		return nil, err
	}

	if po.Status != model.Draft {
		tx.Rollback()
		h.logger.Error("Cannot add items to non-draft purchase order", "status", po.Status)
		return nil, gorm.ErrInvalidData
	}

	// Get product price
	product, err := h.ProductRepo.Search(tx, map[string]interface{}{
		"product_id": req.ProductId,
	}, "")
	if err != nil {
		tx.Rollback()
		h.logger.Error("Product not found", "product_id", req.ProductId, "error", err)
		return nil, err
	}

	// Create item
	item := &model.PurchaseOrderItem{
		PurchaseOrderItemId: uuid.New(),
		PurchaseOrderId:     req.PurchaseOrderId,
		ProductId:           req.ProductId,
		Quantity:            req.Quantity,
		Price:               product.CostPrice,
	}

	if err := h.POItemRepo.Create(tx, item); err != nil {
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

	h.logger.Info("Purchase order item created", "item_id", item.PurchaseOrderItemId)
	return item, nil
}
