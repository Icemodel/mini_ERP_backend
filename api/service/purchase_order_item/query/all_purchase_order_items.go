package query

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/model"

	"gorm.io/gorm"
)

type AllPurchaseOrderItems struct {
	logger   *slog.Logger
	db       *gorm.DB
	ItemRepo repository.PurchaseOrderItem
}

type AllPurchaseOrderItemsRequest struct{}

type AllPurchaseOrderItemsResult struct {
	PurchaseOrderItems []*model.PurchaseOrderItem `json:"purchase_order_items"`
}

func NewAllPurchaseOrderItems(
	logger *slog.Logger,
	db *gorm.DB,
	itemRepo repository.PurchaseOrderItem,
) *AllPurchaseOrderItems {
	return &AllPurchaseOrderItems{
		logger:   logger,
		db:       db,
		ItemRepo: itemRepo,
	}
}

func (h *AllPurchaseOrderItems) Handle(ctx context.Context, req *AllPurchaseOrderItemsRequest) (*AllPurchaseOrderItemsResult, error) {
	items, err := h.ItemRepo.Searches(h.db, map[string]interface{}{}, "")
	if err != nil {
		return nil, err
	}

	return &AllPurchaseOrderItemsResult{PurchaseOrderItems: items}, nil
}