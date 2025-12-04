package query

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PurchaseOrderItems struct {
	logger *slog.Logger
	db     *gorm.DB
	POItemRepo repository.PurchaseOrderItem
}

type PurchaseOrderItemsRequest struct {
	PurchaseOrderId uuid.UUID
}

func NewPurchaseOrderItems(
	logger *slog.Logger,
	db *gorm.DB,
	POItemRepo repository.PurchaseOrderItem,
) *PurchaseOrderItems {
	return &PurchaseOrderItems{
		logger: logger,
		db:     db,
		POItemRepo: POItemRepo,
	}
}

func (h *PurchaseOrderItems) Handle(ctx context.Context, req *PurchaseOrderItemsRequest) (interface{}, error) {
	PO_id := map[string]interface{}{
		"purchase_order_id": req.PurchaseOrderId,
	}
	items, err := h.POItemRepo.Searches(h.db, PO_id, "")
	if err != nil {
		return nil, err
	}

	return items, nil
}
