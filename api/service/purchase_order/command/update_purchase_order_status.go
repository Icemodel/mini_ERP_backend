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
	POItemRepo repository.PurchaseOrderItem
	StockRepo repository.StockTransaction
}

type UpdatePOStatusRequest struct {
	PurchaseOrderId uuid.UUID                 `json:"-" swaggerignore:"true"`
	Status          model.PurchaseOrderStatus `json:"status"`
}

type UpdatePOStatusResult struct {
	PurchaseOrderId uuid.UUID                 `json:"purchase_order_id"`
	Status          model.PurchaseOrderStatus `json:"status"`
}

func NewUpdatePOStatus(
	logger *slog.Logger,
	db *gorm.DB,
	poRepo repository.PurchaseOrder,
	poItemRepo repository.PurchaseOrderItem,
	stockRepo repository.StockTransaction,
) *UpdatePOStatus {
	return &UpdatePOStatus{
		logger:    logger,
		db:        db,
		PORepo:    poRepo,
		POItemRepo: poItemRepo,
		StockRepo: stockRepo,
	}
}

func (h *UpdatePOStatus) Handle(ctx context.Context, req *UpdatePOStatusRequest) (*UpdatePOStatusResult, error) {

	if req.PurchaseOrderId == uuid.Nil {
		return nil, errors.New("purchase_order_id is required")
	}

	if req.Status != model.Confirmed && req.Status != model.Received && req.Status != model.Cancelled {
		return nil, errors.New("invalid status")
	}

	// Begin transaction
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// If changing to RECEIVED, validate that items exist
	purchase_order_id := map[string]interface{}{
		"purchase_order_id": req.PurchaseOrderId,
	}
	items, err := h.POItemRepo.Searches(tx, purchase_order_id, "")
	if req.Status == model.Received {
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

	// Get PurchaseOrder to access CreatedBy
	po, err := h.PORepo.Search(tx, purchase_order_id, "")
	if err != nil {
		tx.Rollback()
		return nil, err
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
			// Fetch latest stock transaction for the product
			latestStockTx, err := h.StockRepo.GetLatestByProduct(tx, item.ProductId)
			if err != nil {
				tx.Rollback()
				return nil, err
			}

			// Calculate new quantity 
			// newQuantity := int64(item.Quantity)
			// if latestStockTx != nil {
			// 	newQuantity = latestStockTx.Quantity + int64(item.Quantity)
			// }

			var newQuantity int64
			if latestStockTx != nil {
				newQuantity = latestStockTx.Quantity + int64(item.Quantity)
			} else {
				newQuantity = int64(item.Quantity)
			}

			received := "Purchase Order Received"

			// Create stock transaction
			stockTx := &model.StockTransaction{
				StockTransactionId: uuid.New(),
				ProductId:          item.ProductId,
				Quantity:           newQuantity,
				Type:               "IN",
				Reason:             &received,
				ReferenceId:        &req.PurchaseOrderId,
				CreatedAt:          time.Now(),
				CreatedBy:          po.CreatedBy,
			}

			if err := h.StockRepo.Create(tx, stockTx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		h.logger.Error("Failed to commit transaction", "error", err)
		return nil, err
	}

	return &UpdatePOStatusResult{
		PurchaseOrderId: req.PurchaseOrderId,
		Status:          req.Status,
	}, nil
}
