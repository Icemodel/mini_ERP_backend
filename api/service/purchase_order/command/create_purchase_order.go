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
	logger *slog.Logger
	db     *gorm.DB
	PORepo repository.PurchaseOrderRepository
}

type CreatePurchaseOrderRequest struct {
	SupplierId uuid.UUID                 `json:"supplier_id" validate:"required"`
	CreatedBy  uuid.UUID                 `json:"created_by" validate:"required"`
	Items      []CreatePurchaseOrderItem `json:"items" validate:"required,min=1,dive"`
}

type CreatePurchaseOrderItem struct {
	ProductId uuid.UUID `json:"product_id" validate:"required"`
	Quantity  uint64    `json:"quantity" validate:"required,min=1"`
	Price     float64   `json:"price" validate:"required,min=0"`
}

func NewCreatePurchaseOrderHandler(
	logger *slog.Logger,
	db *gorm.DB,
	poRepo repository.PurchaseOrderRepository,
) *CreatePurchaseOrder {
	return &CreatePurchaseOrder{
		logger: logger,
		db:     db,
		PORepo: poRepo,
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

	// Calculate total amount from items
	var totalAmount uint64
	for _, it := range req.Items {
		totalAmount += uint64(it.Price * float64(it.Quantity))
	}

	// Create Purchase Order
	po := &model.PurchaseOrder{
		PurchaseOrderId: uuid.New(),
		SupplierId:      req.SupplierId,
		Status:          model.Draft,
		TotalAmount:     totalAmount,
		CreatedAt:       time.Now(),
		CreatedBy:       req.CreatedBy,
	}

	if err := h.PORepo.Create(tx, po); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create Purchase Order Items
	for _, it := range req.Items {
		item := &model.PurchaseOrderItem{
			PurchaseOrderItemId: uuid.New(),
			PurchaseOrderId:     po.PurchaseOrderId,
			ProductId:           it.ProductId,
			Quantity:            it.Quantity,
			Price:               it.Price,
		}
		if err := h.PORepo.CreateItem(tx, item); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		h.logger.Error("Failed to commit transaction", "error", err)
		return nil, err
	}

	// Fetch complete PO with relations
	completePO, err := h.PORepo.FindById(h.db, po.PurchaseOrderId)
	if err != nil {
		return nil, err
	}

	h.logger.Info("Purchase order created successfully", "po_id", po.PurchaseOrderId)
	return completePO, nil
}