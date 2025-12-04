package query

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/model"

	"gorm.io/gorm"
)

type AllPurchaseOrders struct {
	logger *slog.Logger
	db     *gorm.DB
	PORepo repository.PurchaseOrder
}

type AllPurchaseOrdersRequest struct {
	Status  *model.PurchaseOrderStatus `json:"status"`
	OrderBy string                     `json:"order_by"`
}

type AllPurchaseOrdersResult struct {
	PurchaseOrders []*model.PurchaseOrder `json:"purchase_orders"`
}

func NewAllPurchaseOrders(
	logger *slog.Logger,
	db *gorm.DB,
	poRepo repository.PurchaseOrder,
) *AllPurchaseOrders {
	return &AllPurchaseOrders{
		logger: logger,
		db:     db,
		PORepo: poRepo,
	}
}

func (h *AllPurchaseOrders) Handle(ctx context.Context, req *AllPurchaseOrdersRequest) (*AllPurchaseOrdersResult, error) {
	conditions := make(map[string]interface{})

	if req.Status != nil {
		conditions["status"] = *req.Status
	}

	orderBy := req.OrderBy
	if orderBy == "" {
		orderBy = "created_at DESC"
	}

	pos, err := h.PORepo.Searches(h.db, conditions, orderBy)
	if err != nil {
		h.logger.Error("Failed to get all purchase orders", "error", err)
		return nil, err
	}

	return &AllPurchaseOrdersResult{PurchaseOrders: pos}, nil
}