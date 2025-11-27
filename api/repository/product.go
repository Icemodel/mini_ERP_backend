package repository

import (
	"log/slog"
	"mini-erp-backend/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductSearchFilters struct {
	Search     string
	CategoryId *uuid.UUID
}

type Product interface {
	// Get
	Search(db *gorm.DB, conditions map[string]interface{}, orderBy string) (*model.Product, error)
	SearchWithFilters(db *gorm.DB, filters ProductSearchFilters, orderBy string) ([]model.Product, error)
	SearchWithFiltersAndPagination(db *gorm.DB, filters ProductSearchFilters, orderBy string, page, pageSize int) ([]model.Product, int64, error)
	ExitedByProductCode(db *gorm.DB, productCode string) (bool, error)
	// Create
	Create(tx *gorm.DB, product *model.Product) error
	// Update
	Update(tx *gorm.DB, product *model.Product) error
	// Delete
	DeleteById(tx *gorm.DB, productId uuid.UUID) error
}

type product struct {
	logger *slog.Logger
}

func NewProduct(logger *slog.Logger) Product {
	return &product{
		logger: logger,
	}
}

func (p product) Search(db *gorm.DB, conditions map[string]interface{}, orderBy string) (*model.Product, error) {
	products := []model.Product{}
	if err := db.Preload("Category").Where(conditions).Order(orderBy).Limit(1).Find(&products).Error; err != nil {
		p.logger.Error("Failed to search product", slog.String("error", err.Error()))
		return nil, err
	} else {
		if len(products) == 0 {
			err := gorm.ErrRecordNotFound
			p.logger.Error("Product not found", slog.String("error", err.Error()))
			return nil, err
		}
	}
	return &products[0], nil
}

func (p product) SearchWithFilters(db *gorm.DB, filters ProductSearchFilters, orderBy string) ([]model.Product, error) {
	products := []model.Product{}
	query := db.Model(&model.Product{})

	if filters.Search != "" {
		searchPattern := "%" + filters.Search + "%"
		query = query.Where("name LIKE ? OR product_code LIKE ?", searchPattern, searchPattern)
	}

	if filters.CategoryId != nil && *filters.CategoryId != uuid.Nil {
		query = query.Where("category_id = ?", *filters.CategoryId)
	}

	if orderBy != "" {
		query = query.Order(orderBy)
	}

	if err := query.Preload("Category").Find(&products).Error; err != nil {
		p.logger.Error("Failed to search products with filters", slog.String("error", err.Error()))
		return nil, err
	}

	return products, nil
}

func (p product) SearchWithFiltersAndPagination(db *gorm.DB, filters ProductSearchFilters, orderBy string, page, pageSize int) ([]model.Product, int64, error) {
	products := []model.Product{}
	var total int64

	query := db.Model(&model.Product{})

	// ค้นหาจาก search (ค้นหาทั้งชื่อและ product_code)
	if filters.Search != "" {
		searchPattern := "%" + filters.Search + "%"
		query = query.Where("name LIKE ? OR product_code LIKE ?", searchPattern, searchPattern)
	}

	// กรองตาม category_id
	if filters.CategoryId != nil && *filters.CategoryId != uuid.Nil {
		query = query.Where("category_id = ?", filters.CategoryId)
	}

	// นับจำนวนทั้งหมด
	if err := query.Count(&total).Error; err != nil {
		p.logger.Error("Failed to count products with filters", slog.String("error", err.Error()))
		return nil, 0, err
	}

	// คำนวณ offset
	offset := (page - 1) * pageSize

	// เรียงลำดับ
	if orderBy != "" {
		query = query.Order(orderBy)
	}

	// ดึงข้อมูลแบบ pagination พร้อม preload Category
	if err := query.Preload("Category").Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		p.logger.Error("Failed to search products with filters and pagination", slog.String("error", err.Error()))
		return nil, 0, err
	}

	return products, total, nil
}

func (p product) ExitedByProductCode(db *gorm.DB, productCode string) (bool, error) {
	var count int64
	if err := db.Model(&model.Product{}).Where("product_code = ?", productCode).Count(&count).Error; err != nil {
		p.logger.Error("Failed to check if product code exists", slog.String("error", err.Error()))
		return false, err
	}
	return count > 0, nil
}

func (p product) Create(tx *gorm.DB, product *model.Product) error {
	if err := tx.Create(product).Error; err != nil {
		p.logger.Error("Failed to create product", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (p product) Update(tx *gorm.DB, product *model.Product) error {
	if err := tx.Save(product).Error; err != nil {
		p.logger.Error("Failed to update product", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (p product) DeleteById(tx *gorm.DB, productId uuid.UUID) error {
	if err := tx.Delete(&model.Product{}, "product_id = ?", productId).Error; err != nil {
		p.logger.Error("Failed to delete product", slog.String("error", err.Error()))
		return err
	}
	return nil
}
