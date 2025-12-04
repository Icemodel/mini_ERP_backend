package repository

import (
	"log/slog"
	"mini-erp-backend/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PurchaseOrder interface {
	Create(tx *gorm.DB, po *model.PurchaseOrder) error
	Update(tx *gorm.DB, po *model.PurchaseOrder) error
	UpdateStatus(tx *gorm.DB, poId uuid.UUID, status model.PurchaseOrderStatus) error
	Search(db *gorm.DB, conditions map[string]interface{}, orderBy string) (*model.PurchaseOrder, error)
	Searches(db *gorm.DB, conditions map[string]interface{}, orderBy string) ([]*model.PurchaseOrder, error)
}

type purchaseOrder struct {
	logger *slog.Logger
}

func NewPurchaseOrder(logger *slog.Logger) PurchaseOrder {
	return &purchaseOrder{
		logger: logger,
	}
}

func (r *purchaseOrder) Create(tx *gorm.DB, po *model.PurchaseOrder) error {
	err := tx.Create(po).Error
	if err != nil {
		r.logger.Error("Failed to create purchase order", "error", err)
		return err
	}
	return nil
}

func (r *purchaseOrder) Update(tx *gorm.DB, po *model.PurchaseOrder) error {
	if err := tx.Model(&model.PurchaseOrder{}).
		Where("purchase_order_id = ?", po.PurchaseOrderId).
		Select("*").
		Omit("created_at", "purchase_order_id").
		Updates(po).Error; err != nil {
		r.logger.Error("Failed to update purchase order", "error", err)
		return err
	}
	return nil
}

func (r *purchaseOrder) UpdateStatus(tx *gorm.DB, poId uuid.UUID, status model.PurchaseOrderStatus) error {
	if err := tx.Model(&model.PurchaseOrder{}).
		Where("purchase_order_id = ?", poId).
		Update("status", status).Error; err != nil {
		r.logger.Error("Failed to update purchase order status", "error", err)
		return err
	}
	return nil
}

func (r *purchaseOrder) Search(db *gorm.DB, conditions map[string]interface{}, orderBy string) (*model.PurchaseOrder, error) {
	pos := []model.PurchaseOrder{}

	if err := db.Preload("Supplier").Preload("PurchaseOrderItem").Where(conditions).Order(orderBy).Limit(1).Find(&pos).Error; err != nil {
		r.logger.Error("Failed to search purchase order", "error", err)
		return nil, err
	} else {
		if len(pos) == 0 {
			err := gorm.ErrRecordNotFound
			r.logger.Error("Purchase order not found", "error", err)
			return nil, err
		}
	}

	return &pos[0], nil
}

func (r *purchaseOrder) Searches(db *gorm.DB, conditions map[string]interface{}, orderBy string) ([]*model.PurchaseOrder, error) {
	pos := []*model.PurchaseOrder{}

	if err := db.Preload("Supplier").Preload("PurchaseOrderItem").Where(conditions).Order(orderBy).Find(&pos).Error; err != nil {
		r.logger.Error("Failed to search purchase orders", "error", err)
		return nil, err
	} else {
		if len(pos) == 0 {
			err := gorm.ErrRecordNotFound
			r.logger.Error("No purchase orders found", "error", err)
			return nil, err
		}
	}

	return pos, nil
}