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

	CreateItem(tx *gorm.DB, item *model.PurchaseOrderItem) error
	DeleteItemsByPurchaseOrderId(tx *gorm.DB, poId uuid.UUID) error
	SearchItemsByPurchaseOrderId(db *gorm.DB, poId uuid.UUID) ([]*model.PurchaseOrderItem, error)
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
		Omit("created_at", "created_by", "purchase_order_id").
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

	query := db.Preload("PurchaseOrderItem").
		Preload("PurchaseOrderItem.Product").
		Preload("Supplier").
		Where(conditions)

	if orderBy != "" {
		query = query.Order(orderBy)
	}

	if err := query.Find(&pos).Error; err != nil {
		r.logger.Error("Failed to search purchase order", "error", err)
		return nil, err
	}

	if len(pos) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &pos[0], nil
}

func (r *purchaseOrder) Searches(db *gorm.DB, conditions map[string]interface{}, orderBy string) ([]*model.PurchaseOrder, error) {
	pos := []*model.PurchaseOrder{}

	query := db.Preload("Supplier").Preload("PurchaseOrderItem").Where(conditions)

	if orderBy != "" {
		query = query.Order(orderBy)
	}

	if err := query.Find(&pos).Error; err != nil {
		r.logger.Error("Failed to search purchase orders", "error", err)
		return nil, err
	}
	return pos, nil
}

func (r *purchaseOrder) CreateItem(tx *gorm.DB, item *model.PurchaseOrderItem) error {
	if err := tx.Create(item).Error; err != nil {
		r.logger.Error("Failed to create purchase order item", "error", err)
		return err
	}
	return nil
}

func (r *purchaseOrder) DeleteItemsByPurchaseOrderId(tx *gorm.DB, poId uuid.UUID) error {
	if err := tx.Where("purchase_order_id = ?", poId).Delete(&model.PurchaseOrderItem{}).Error; err != nil {
		r.logger.Error("Failed to delete purchase order items", "error", err)
		return err
	}
	return nil
}

func (r *purchaseOrder) SearchItemsByPurchaseOrderId(db *gorm.DB, poId uuid.UUID) ([]*model.PurchaseOrderItem, error) {
	items := []*model.PurchaseOrderItem{}
	if err := db.Where("purchase_order_id = ?", poId).
		Find(&items).Error; err != nil {
		r.logger.Error("Failed to search purchase order items", "error", err)
		return nil, err
	}
	return items, nil
}
