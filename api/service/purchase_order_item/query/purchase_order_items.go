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
	PORepo repository.PurchaseOrder
}

type PurchaseOrderItemsRequest struct {
	PurchaseOrderId uuid.UUID
}

func NewPurchaseOrderItems(
	logger *slog.Logger,
	db *gorm.DB,
	poRepo repository.PurchaseOrder,
) *PurchaseOrderItems {
	return &PurchaseOrderItems{
		logger: logger,
		db:     db,
		PORepo: poRepo,
	}
}

func (h *PurchaseOrderItems) Handle(ctx context.Context, req *PurchaseOrderItemsRequest) (interface{}, error) {
	items, err := h.PORepo.SearchItemsByPurchaseOrderId(h.db, req.PurchaseOrderId)
	if err != nil {
		return nil, err
	}

	return items, nil
}
