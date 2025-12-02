package query

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetPurchaseOrder struct {
	logger *slog.Logger
	db     *gorm.DB
	PORepo repository.PurchaseOrder
}

type GetPurchaseOrderRequest struct {
	PurchaseOrderId uuid.UUID `json:"purchase_order_id" validate:"required"`
}

type GetPurchaseOrderResult struct {
	PurchaseOrder *model.PurchaseOrder `json:"purchase_order"`
}

func NewGetPurchaseOrder(
	logger *slog.Logger,
	db *gorm.DB,
	poRepo repository.PurchaseOrder,
) *GetPurchaseOrder {
	return &GetPurchaseOrder{
		logger: logger,
		db:     db,
		PORepo: poRepo,
	}
}

func (h *GetPurchaseOrder) Handle(ctx context.Context, req *GetPurchaseOrderRequest) (*GetPurchaseOrderResult, error) {
	po, err := h.PORepo.Search(h.db, map[string]interface{}{
		"purchase_order_id": req.PurchaseOrderId,
	}, "")
	if err != nil {
		h.logger.Error("Failed to get purchase order", "po_id", req.PurchaseOrderId, "error", err)
		return nil, err
	}

	return &GetPurchaseOrderResult{PurchaseOrder: po}, nil
}
