package query

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PurchaseOrder struct {
	logger *slog.Logger
	db     *gorm.DB
	PORepo repository.PurchaseOrder
}

type PurchaseOrderRequest struct {
	PurchaseOrderId uuid.UUID `json:"purchase_order_id" validate:"required"`
}

type PurchaseOrderResult struct {
	PurchaseOrder *model.PurchaseOrder `json:"purchase_order"`
}

func NewPurchaseOrder(
	logger *slog.Logger,
	db *gorm.DB,
	poRepo repository.PurchaseOrder,
) *PurchaseOrder {
	return &PurchaseOrder{
		logger: logger,
		db:     db,
		PORepo: poRepo,
	}
}

func (h *PurchaseOrder) Handle(ctx context.Context, req *PurchaseOrderRequest) (*PurchaseOrderResult, error) {
	condition := map[string]interface{}{
		"purchase_order_id": req.PurchaseOrderId,
	}
	po, err := h.PORepo.Search(h.db, condition, "")
	if err != nil {
		h.logger.Error("Failed to get purchase order", "po_id", req.PurchaseOrderId, "error", err)
		return nil, err
	}

	return &PurchaseOrderResult{PurchaseOrder: po}, nil
}