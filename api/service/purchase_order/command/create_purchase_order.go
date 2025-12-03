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

type CreatePurchaseOrder struct {
	logger      *slog.Logger
	db          *gorm.DB
	PORepo      repository.PurchaseOrder
	ProductRepo repository.Product
}

type CreatePurchaseOrderRequest struct {
	SupplierId uuid.UUID `json:"supplier_id"`
	CreatedBy  string    `json:"created_by"`
}

type CreatePurchaseOrderResult struct {
	PurchaseOrder model.PurchaseOrder `json:"purchase_order"`
}

func NewCreatePurchaseOrder(
	logger *slog.Logger,
	db *gorm.DB,
	poRepo repository.PurchaseOrder,
	productRepo repository.Product,
) *CreatePurchaseOrder {
	return &CreatePurchaseOrder{
		logger:      logger,
		db:          db,
		PORepo:      poRepo,
		ProductRepo: productRepo,
	}
}

func (h *CreatePurchaseOrder) Handle(ctx context.Context, req *CreatePurchaseOrderRequest) (*CreatePurchaseOrderResult, error) {

	if req.SupplierId == uuid.Nil {
		return nil, errors.New("supplier_id is required")
	}

	if req.CreatedBy == "" {
		return nil, errors.New("created_by is required")
	}

	// Begin transaction
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create Purchase Order
	po := &model.PurchaseOrder{
		PurchaseOrderId: uuid.New(),
		SupplierId:      req.SupplierId,
		Status:          model.Draft,
		CreatedAt:       time.Now(),
		CreatedBy:       req.CreatedBy,
	}

	if err := h.PORepo.Create(tx, po); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		h.logger.Error("Failed to commit transaction", "error", err)
		return nil, err
	}

	return &CreatePurchaseOrderResult{
		PurchaseOrder: *po,
	}, nil
}
