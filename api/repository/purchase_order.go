package repository

import (
	"log/slog"
	"mini-erp-backend/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PurchaseOrderRepository interface {
	Create(tx *gorm.DB, po *model.PurchaseOrder) error
	Update(tx *gorm.DB, po *model.PurchaseOrder) error
	UpdateStatus(tx *gorm.DB, poId uuid.UUID, status model.PurchaseOrderStatus) error
	FindById(db *gorm.DB, poId uuid.UUID) (*model.PurchaseOrder, error)
	FindAll(db *gorm.DB, conditions map[string]interface{}, orderBy string) ([]*model.PurchaseOrder, error)
	
	CreateItem(tx *gorm.DB, item *model.PurchaseOrderItem) error
	DeleteItemsByPOId(tx *gorm.DB, poId uuid.UUID) error
	FindItemsByPOId(db *gorm.DB, poId uuid.UUID) ([]*model.PurchaseOrderItem, error)
}

type purchaseOrderRepository struct {
	logger *slog.Logger
}

func NewPurchaseOrderRepository(logger *slog.Logger) PurchaseOrderRepository {
	return &purchaseOrderRepository{
		logger: logger,
	}
}

func (r *purchaseOrderRepository) Create(tx *gorm.DB, po *model.PurchaseOrder) error {
	err := tx.Create(po).Error
	if err != nil {
		r.logger.Error("Failed to create purchase order", "error", err)
		return err
	}
	return nil
}

func (r *purchaseOrderRepository) Update(tx *gorm.DB, po *model.PurchaseOrder) error {
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

func (r *purchaseOrderRepository) UpdateStatus(tx *gorm.DB, poId uuid.UUID, status model.PurchaseOrderStatus) error {
	if err := tx.Model(&model.PurchaseOrder{}).
		Where("purchase_order_id = ?", poId).
		Update("status", status).Error; err != nil {
		r.logger.Error("Failed to update purchase order status", "error", err)
		return err
	}
	return nil
}

func (r *purchaseOrderRepository) FindById(db *gorm.DB, poId uuid.UUID) (*model.PurchaseOrder, error) {
	var po model.PurchaseOrder
	if err := db.Preload("PurchaseOrderItem").
		Preload("PurchaseOrderItem.Product").
		Preload("Supplier").
		Where("purchase_order_id = ?", poId).
		First(&po).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Error("Purchase order not found", "po_id", poId)
		} else {
			r.logger.Error("Failed to find purchase order", "error", err)
		}
		return nil, err
	}
	return &po, nil
}

func (r *purchaseOrderRepository) FindAll(db *gorm.DB, conditions map[string]interface{}, orderBy string) ([]*model.PurchaseOrder, error) {
	var pos []*model.PurchaseOrder
	query := db.Preload("Supplier").Preload("PurchaseOrderItem")
	
	if len(conditions) > 0 {
		query = query.Where(conditions)
	}
	
	if orderBy != "" {
		query = query.Order(orderBy)
	}
	
	if err := query.Find(&pos).Error; err != nil {
		r.logger.Error("Failed to find purchase orders", "error", err)
		return nil, err
	}
	return pos, nil
}

func (r *purchaseOrderRepository) CreateItem(tx *gorm.DB, item *model.PurchaseOrderItem) error {
	if err := tx.Create(item).Error; err != nil {
		r.logger.Error("Failed to create purchase order item", "error", err)
		return err
	}
	return nil
}

func (r *purchaseOrderRepository) DeleteItemsByPOId(tx *gorm.DB, poId uuid.UUID) error {
	if err := tx.Where("purchase_order_id = ?", poId).Delete(&model.PurchaseOrderItem{}).Error; err != nil {
		r.logger.Error("Failed to delete purchase order items", "error", err)
		return err
	}
	return nil
}

func (r *purchaseOrderRepository) FindItemsByPOId(db *gorm.DB, poId uuid.UUID) ([]*model.PurchaseOrderItem, error) {
	var items []*model.PurchaseOrderItem
	if err := db.Preload("Product").
		Where("purchase_order_id = ?", poId).
		Find(&items).Error; err != nil {
		r.logger.Error("Failed to find purchase order items", "error", err)
		return nil, err
	}
	return items, nil
}
