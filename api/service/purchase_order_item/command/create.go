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
	PurchaseOrderId uuid.UUID `json:"purchase_order_id"`
	ProductId       uuid.UUID `json:"product_id"`
	Quantity        uint64    `json:"quantity"`
}

type CreatePurchaseOrderItemResult struct {
	PurchaseOrderItem model.PurchaseOrderItem `json:"purchase_order_item"`
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

func (h *CreatePurchaseOrderItem) Handle(ctx context.Context, req *CreatePurchaseOrderItemRequest) (*CreatePurchaseOrderItemResult, error) {

	if req.PurchaseOrderId == uuid.Nil {
		h.logger.Error("purchase_order_id is required")
		return nil, gorm.ErrInvalidData
	}

	if req.ProductId == uuid.Nil {
		h.logger.Error("product_id is required")
		return nil, gorm.ErrInvalidData
	}

	if req.Quantity < 1 {
		h.logger.Error("quantity must be at least 1")
		return nil, gorm.ErrInvalidData
	}

	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Verify PO exists and is DRAFT
	purchase_order_id := map[string]interface{}{
		"purchase_order_id": req.PurchaseOrderId,
	}
	po, err := h.PORepo.Search(tx, purchase_order_id, "")
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
	condition := map[string]interface{}{
		"product_id": req.ProductId,
	}
	product, err := h.ProductRepo.Search(tx, condition, "")
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
	

	if err := tx.Commit().Error; err != nil {
		h.logger.Error("Failed to commit transaction", "error", err)
		return nil, err
	}

	return &CreatePurchaseOrderItemResult{
		PurchaseOrderItem: *item,
	}, nil
}
