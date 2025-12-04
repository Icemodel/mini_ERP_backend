package repository

import (
	"log/slog"
	"mini-erp-backend/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PurchaseOrderItem interface {
	Create(tx *gorm.DB, item *model.PurchaseOrderItem) error
	Delete(tx *gorm.DB, itemId uuid.UUID) error
	Search(db *gorm.DB, conditions map[string]interface{}, orderBy string) (*model.PurchaseOrderItem, error)
	Searches(db *gorm.DB, conditions map[string]interface{}, orderBy string) ([]*model.PurchaseOrderItem, error)
	Update(tx *gorm.DB, itemId uuid.UUID, item *model.PurchaseOrderItem) error
}

type purchaseOrderItem struct {
	logger *slog.Logger
}

func NewPurchaseOrderItem(logger *slog.Logger) PurchaseOrderItem {
	return &purchaseOrderItem{
		logger: logger,
	}
}

func (r *purchaseOrderItem) Create(tx *gorm.DB, item *model.PurchaseOrderItem) error {
	err := tx.Create(item).Error
	if err != nil {
		r.logger.Error("Failed to create purchase order item", "error", err)
		return err
	}
	return nil
}

func (r *purchaseOrderItem) Delete(tx *gorm.DB, itemId uuid.UUID) error {
	if err := tx.Delete(&model.PurchaseOrderItem{}, "purchase_order_item_id = ?", itemId).Error; err != nil {
		r.logger.Error("Failed to delete purchase order item", "error", err)
		return err
	}
	return nil
}

func (r *purchaseOrderItem) Search(db *gorm.DB, conditions map[string]interface{}, orderBy string) (*model.PurchaseOrderItem, error) {
	items := []model.PurchaseOrderItem{}
	if err := db.Preload("Product").Where(conditions).Order(orderBy).Limit(1).Find(&items).Error; err != nil {
		r.logger.Error("Failed to get purchase order item", "error", err)
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}
	return &items[0], nil
}

func (r *purchaseOrderItem) Searches(db *gorm.DB, conditions map[string]interface{}, orderBy string) ([]*model.PurchaseOrderItem, error) {
	items := []*model.PurchaseOrderItem{}
	if err := db.Preload("Product").Where(conditions).Order(orderBy).Find(&items).Error; err != nil {
		r.logger.Error("Failed to get purchase order items", "error", err)
		return nil, err
	}
	return items, nil
}

func (r *purchaseOrderItem) Update(tx *gorm.DB, itemId uuid.UUID, item *model.PurchaseOrderItem) error {
	if err := tx.Model(&model.PurchaseOrderItem{}).Where("purchase_order_item_id = ?", itemId).Select("*").Omit("created_at").Updates(item).Error; err != nil {
		r.logger.Error("Failed to update purchase order item", "error", err)
		return err
	}
	return nil
}
