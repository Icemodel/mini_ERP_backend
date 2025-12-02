package command

import (
	"context"
	"errors"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdatePOStatus struct {
	logger    *slog.Logger
	db        *gorm.DB
	PORepo    repository.PurchaseOrder
	StockRepo repository.StockTransaction
}

type UpdatePOStatusRequest struct {
	PurchaseOrderId uuid.UUID                 
	Status          model.PurchaseOrderStatus `json:"status" validate:"required"`
	CreatedBy       uuid.UUID                 `json:"created_by" validate:"required"`
}

func NewUpdatePOStatus(
	logger *slog.Logger,
	db *gorm.DB,
	poRepo repository.PurchaseOrder,
	stockRepo repository.StockTransaction,
) *UpdatePOStatus {
	return &UpdatePOStatus{
		logger:    logger,
		db:        db,
		PORepo:    poRepo,
		StockRepo: stockRepo,
	}
}

func (h *UpdatePOStatus) Handle(ctx context.Context, req *UpdatePOStatusRequest) (interface{}, error) {
	// Begin transaction
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// If changing to RECEIVED, validate that items exist
	var items []*model.PurchaseOrderItem
	var err error
	if req.Status == model.Received {
		items, err = h.PORepo.SearchItemsByPurchaseOrderId(tx, req.PurchaseOrderId)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if len(items) == 0 {
			tx.Rollback()
			h.logger.Error("Cannot receive purchase order without items", "po_id", req.PurchaseOrderId)
			return nil, errors.New("cannot receive purchase order without items")
		}
	}

	// Update status
	if err := h.PORepo.UpdateStatus(tx, req.PurchaseOrderId, req.Status); err != nil {
		tx.Rollback()
		return nil, err
	}

	// If status is RECEIVED, create Stock IN transactions
	if req.Status == model.Received {
		// Create Stock IN transactions for each item
		for _, item := range items {
			// หา transaction ล่าสุด สำหรับ product_id
			latestStockTx, err := h.StockRepo.GetLatestByProduct(tx, item.ProductId)
			if err != nil {
				tx.Rollback()
				return nil, err
			}

			// คำนวณ quantity ใหม่: PO item quantity + quantity ของ transaction ล่าสุด
			newQuantity := int64(item.Quantity)
			if latestStockTx != nil {
				newQuantity = latestStockTx.Quantity + int64(item.Quantity)
			}

			// สร้าง transaction ใหม่
			stockTx := &model.StockTransaction{
				StockTransactionId: uuid.New(),
				ProductId:          item.ProductId,
				Quantity:           newQuantity,
				Type:               "IN",
				Reason:             stringPtr("Purchase Order Received"),
				ReferenceId:        &req.PurchaseOrderId,
				CreatedAt:          time.Now(),
				CreatedBy:          req.CreatedBy.String(),
			}

			if err := h.StockRepo.Create(tx, stockTx); err != nil {
				tx.Rollback()
				return nil, err
			}

			if latestStockTx != nil {
				h.logger.Info("Stock transaction created with cumulative quantity",
					"product_id", item.ProductId,
					"po_item_quantity", item.Quantity,
					"previous_quantity", latestStockTx.Quantity,
					"new_quantity", newQuantity,
					"po_id", req.PurchaseOrderId)
			} else {
				h.logger.Info("Stock transaction created (first)",
					"product_id", item.ProductId,
					"quantity", newQuantity,
					"po_id", req.PurchaseOrderId)
			}
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		h.logger.Error("Failed to commit transaction", "error", err)
		return nil, err
	}

	h.logger.Info("Purchase order status updated", "po_id", req.PurchaseOrderId, "status", req.Status)
	return map[string]interface{}{
		"purchase_order_id": req.PurchaseOrderId,
		"status":            req.Status,
	}, nil
}

func stringPtr(s string) *string {
	return &s
}
