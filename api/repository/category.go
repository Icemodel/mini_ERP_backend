package repository

import (
	"log/slog"
	"mini-erp-backend/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategorySearchFilters struct {
	Search string // ค้นหาทั้งชื่อและคำอธิบาย
}

type Category interface {
	// Get
	Search(db *gorm.DB, conditions map[string]interface{}, orderBy string) (*model.Category, error)
	// Searches(db *gorm.DB, conditions map[string]interface{}, orderBy string) ([]model.Category, error)
	// SearchesWithPagination(db *gorm.DB, conditions map[string]interface{}, orderBy string, offset, limit int) ([]model.Category, int64, error)
	SearchWithFilters(db *gorm.DB, filters CategorySearchFilters, orderBy string) ([]model.Category, error)
	SearchWithFiltersAndPagination(db *gorm.DB, filters CategorySearchFilters, orderBy string, page, pageSize int) ([]model.Category, int64, error)
	ExitedByName(db *gorm.DB, name string) (bool, error)
	ExitedByNameExcludeId(db *gorm.DB, name string, categoryId uuid.UUID) (bool, error)
	// Create
	Create(tx *gorm.DB, category *model.Category) error
	// Update
	Update(tx *gorm.DB, category *model.Category) error
	// Delete
	DeleteById(tx *gorm.DB, categoryId uuid.UUID) error
}

type category struct {
	logger *slog.Logger
}

func NewCategory(logger *slog.Logger) Category {
	return &category{
		logger: logger,
	}
}

func (c category) Search(db *gorm.DB, conditions map[string]interface{}, orderBy string) (*model.Category, error) {
	category := []model.Category{}

	if err := db.Where(conditions).Order(orderBy).Limit(1).Find(&category).Error; err != nil {
		c.logger.Error("Failed to search category", slog.String("error", err.Error()))
		return nil, err
	} else {
		if len(category) == 0 {
			err = gorm.ErrRecordNotFound
			c.logger.Error("Category not found", slog.String("error", err.Error()))
			return nil, err
		}
	}

	return &category[0], nil
}

func (c category) ExitedByName(db *gorm.DB, name string) (bool, error) {
	var count int64
	if err := db.Model(&model.Category{}).Where("name = ?", name).Count(&count).Error; err != nil {
		c.logger.Error("Failed to check if category exists by name", slog.String("error", err.Error()))
		return false, err
	}
	return count > 0, nil
}

func (c category) ExitedByNameExcludeId(db *gorm.DB, name string, categoryId uuid.UUID) (bool, error) {
	var count int64
	if err := db.Model(&model.Category{}).Where("name = ? AND category_id != ?", name, categoryId).Count(&count).Error; err != nil {
		c.logger.Error("Failed to check if category exists by name excluding id", slog.String("error", err.Error()))
		return false, err
	}
	return count > 0, nil
}

func (c category) Create(tx *gorm.DB, category *model.Category) error {
	if err := tx.Create(category).Error; err != nil {
		c.logger.Error("Failed to create category", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (c category) Update(tx *gorm.DB, category *model.Category) error {
	if err := tx.Save(category).Error; err != nil {
		c.logger.Error("Failed to update category", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (c category) DeleteById(tx *gorm.DB, categoryId uuid.UUID) error {
	if err := tx.Delete(&model.Category{}, "category_id = ?", categoryId).Error; err != nil {
		c.logger.Error("Failed to delete category", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (c category) SearchWithFilters(db *gorm.DB, filters CategorySearchFilters, orderBy string) ([]model.Category, error) {
	categories := []model.Category{}
	query := db.Model(&model.Category{})

	// ค้นหาจาก search (ค้นหาทั้งชื่อและ description)
	if filters.Search != "" {
		searchPattern := "%" + filters.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	// เรียงลำดับ
	if orderBy != "" {
		query = query.Order(orderBy)
	}

	if err := query.Find(&categories).Error; err != nil {
		c.logger.Error("Failed to search categories with filters", slog.String("error", err.Error()))
		return nil, err
	}

	return categories, nil
}

func (c category) SearchWithFiltersAndPagination(db *gorm.DB, filters CategorySearchFilters, orderBy string, page, pageSize int) ([]model.Category, int64, error) {
	categories := []model.Category{}
	var total int64

	query := db.Model(&model.Category{})

	// ค้นหาจาก search (ค้นหาทั้งชื่อและ description)
	if filters.Search != "" {
		searchPattern := "%" + filters.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	// นับจำนวนทั้งหมด
	if err := query.Count(&total).Error; err != nil {
		c.logger.Error("Failed to count categories with filters", slog.String("error", err.Error()))
		return nil, 0, err
	}

	// คำนวณ offset
	offset := (page - 1) * pageSize

	// เรียงลำดับ
	if orderBy != "" {
		query = query.Order(orderBy)
	}

	// ดึงข้อมูลแบบ pagination
	if err := query.Offset(offset).Limit(pageSize).Find(&categories).Error; err != nil {
		c.logger.Error("Failed to search categories with filters and pagination", slog.String("error", err.Error()))
		return nil, 0, err
	}

	return categories, total, nil
}
