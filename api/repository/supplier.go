package repository

import (
	"log/slog"
	"mini-erp-backend/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Supplier interface {
	Create(tx *gorm.DB, supplier *model.Supplier) error
	UpdateBySupplierId(tx *gorm.DB, supplierId uuid.UUID, supplier *model.Supplier) error
	Delete(tx *gorm.DB, supplier *model.Supplier) error
	Search(db *gorm.DB, conditions map[string]interface{}, orderBy string) (*model.Supplier, error)
	Searches(db *gorm.DB, conditions map[string]interface{}, orderBy string) ([]*model.Supplier, error)
}

type supplier struct {
	logger *slog.Logger
}

func NewSupplier(logger *slog.Logger) Supplier {
	return &supplier{
		logger: logger,
	}
}

func (r *supplier) Create(tx *gorm.DB, supplier *model.Supplier) error {
	err := tx.Create(supplier).Error
	if err != nil {
		r.logger.Error("Failed to create supplier", "error", err)
		return err
	}

	return nil
}

func (r *supplier) UpdateBySupplierId(tx *gorm.DB, supplierId uuid.UUID, supplier *model.Supplier) error {
	if err := tx.Model(&model.Supplier{}).Where("supplier_id = ?", supplierId).Select("*").Omit("created_at").Updates(supplier).Error; err != nil {
		r.logger.Error("Failed to update supplier", "error", err)
		return err
	}
	return nil
}

func (r *supplier) Delete(tx *gorm.DB, supplier *model.Supplier) error {
	if err := tx.Delete(supplier).Error; err != nil {
		r.logger.Error("Failed to delete supplier", "error", err)
		return err
	}
	return nil
}

func (r *supplier) Search(db *gorm.DB, conditions map[string]interface{}, orderBy string) (*model.Supplier, error) {
	suppliers := []model.Supplier{}

	if err := db.Where(conditions).Order(orderBy).Limit(1).Find(&suppliers).Error; err != nil {
		r.logger.Error("Failed to get supplier", "error", err)
		return nil, err
	} else {
		if len(suppliers) == 0 {
			err = gorm.ErrRecordNotFound
			r.logger.Error("Supplier not found", "error", err)
			return nil, err
		}
	}

	return &suppliers[0], nil
}

func (r *supplier) Searches(db *gorm.DB, conditions map[string]interface{}, orderBy string) ([]*model.Supplier, error) {
	suppliers := []*model.Supplier{}
	if err := db.Where(conditions).Order(orderBy).Find(&suppliers).Error; err != nil {
		r.logger.Error("Failed to get all suppliers", "error", err)
		return nil, err
	}
	return suppliers, nil
}