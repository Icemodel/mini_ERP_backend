package command

import (
	"context"
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
	SupplierId uuid.UUID `json:"supplier_id" validate:"required"`
	CreatedBy  string    `json:"created_by" validate:"required"`
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

func (h *CreatePurchaseOrder) Handle(ctx context.Context, req *CreatePurchaseOrderRequest) (interface{}, error) {
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

	h.logger.Info("Purchase order created successfully", "po_id", po.PurchaseOrderId)
	return po, nil
}