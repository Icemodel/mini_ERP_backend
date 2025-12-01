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
	SupplierId uuid.UUID       `json:"supplier_id" validate:"required"`
	Items      []POItemRequest `json:"items" validate:"required,min=1"`
	CreatedBy  string          `json:"created_by" validate:"required"`
}

type POItemRequest struct {
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
	// Calculate total amount
	var totalAmount uint64
	for _, item := range req.Items {
		totalAmount += uint64(item.Price * float64(item.Quantity))
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
		TotalAmount:     totalAmount,
		CreatedAt:       time.Now(),
		CreatedBy:       req.CreatedBy,
	}

	if err := h.PORepo.Create(tx, po); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create Purchase Order Items
	for _, itemReq := range req.Items {
		item := &model.PurchaseOrderItem{
			PurchaseOrderItemId: uuid.New(),
			PurchaseOrderId:     po.PurchaseOrderId,
			ProductId:           itemReq.ProductId,
			Quantity:            itemReq.Quantity,
			Price:               itemReq.Price,
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
