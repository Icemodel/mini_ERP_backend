package query

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PurchaseOrderItem struct {
	logger     *slog.Logger
	db         *gorm.DB
	POItemRepo repository.PurchaseOrderItem
}

type PurchaseOrderItemRequest struct {
	PurchaseOrderItemId uuid.UUID `json:"purchase_order_item_id"`
}

type PurchaseOrderItemResult struct {
	PurchaseOrderItem model.PurchaseOrderItem `json:"purchase_order_item"`
}

func NewPurchaseOrderItem(
	logger *slog.Logger,
	db *gorm.DB,
	poItemRepo repository.PurchaseOrderItem,
) *PurchaseOrderItem {
	return &PurchaseOrderItem{
		logger:     logger,
		db:         db,
		POItemRepo: poItemRepo,
	}
}

func (h *PurchaseOrderItem) Handle(ctx context.Context, req *PurchaseOrderItemRequest) (*PurchaseOrderItemResult, error) {
	item_id := map[string]interface{}{
		"purchase_order_item_id": req.PurchaseOrderItemId,
	}
	item, err := h.POItemRepo.Search(h.db, item_id, "")
	if err != nil {
		return nil, err
	}

	return &PurchaseOrderItemResult{PurchaseOrderItem: *item}, nil
}