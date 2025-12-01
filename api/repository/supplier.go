package repository

import (
	"log/slog"
	"mini-erp-backend/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)
type SupplierRepository interface {
	Create(tx *gorm.DB, supplier *model.Supplier) error
	UpdateBySupplierId(tx *gorm.DB, supplierId uuid.UUID, supplier *model.Supplier) error
	Delete(tx *gorm.DB, supplier *model.Supplier) error
	Search(db *gorm.DB, conditions map[string]interface{}, orderBy string) (*model.Supplier, error)
	Searches(db *gorm.DB, conditions map[string]interface{}, orderBy string) ([]*model.Supplier, error)
}

type supplierRepository struct {
	logger *slog.Logger
}

func NewSupplierRepository(logger *slog.Logger) SupplierRepository {
	return &supplierRepository{
		logger: logger,
	}
}

func (r *supplierRepository) Create(tx *gorm.DB, supplier *model.Supplier) error {
	err := tx.Create(supplier).Error
	if err != nil {
		r.logger.Error("Failed to create supplier", "error", err)
		return err
	}

	return nil
}

func (r *supplierRepository) UpdateBySupplierId(tx *gorm.DB, supplierId uuid.UUID, supplier *model.Supplier) error {
	if err := tx.Model(&model.Supplier{}).Where("supplier_id = ?", supplierId).Select("*").Omit("created_at").Updates(supplier).Error; err != nil {
		r.logger.Error("Failed to update supplier", "error", err)
		return err
	}
	return nil
}

func (r *supplierRepository) Delete(tx *gorm.DB, supplier *model.Supplier) error {
	if err := tx.Delete(supplier).Error; err != nil {
		r.logger.Error("Failed to delete supplier", "error", err)
		return err
	}
	return nil
}

func (r *supplierRepository) Search(db *gorm.DB, conditions map[string]interface{}, orderBy string) (*model.Supplier, error) {
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

func (r *supplierRepository) Searches(db *gorm.DB, conditions map[string]interface{}, orderBy string) ([]*model.Supplier, error) {
	suppliers := []*model.Supplier{}
	if err := db.Where(conditions).Order(orderBy).Find(&suppliers).Error; err != nil {
		r.logger.Error("Failed to get all suppliers", "error", err)
		return nil, err
	}
	return suppliers, nil
}

