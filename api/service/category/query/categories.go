package query

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/model"

	"gorm.io/gorm"
)

type Categories struct {
	logger       *slog.Logger
	db           *gorm.DB
	categoryRepo repository.Category
}

type CategoriesRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Search   string `json:"search"` // ค้นหาทั้งชื่อและคำอธิบาย
}

type CategoriesResult struct {
	Categories []model.Category `json:"categories"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}

func NewCategories(logger *slog.Logger, db *gorm.DB, categoryRepo repository.Category) *Categories {
	return &Categories{
		logger:       logger,
		db:           db,
		categoryRepo: categoryRepo,
	}
}

func (c *Categories) Handle(ctx context.Context, request CategoriesRequest) (*CategoriesResult, error) {
	// สร้าง filters จาก request
	filters := repository.CategorySearchFilters{
		Search: request.Search,
	}

	// ถ้าไม่มี pagination ให้ดึงทั้งหมด
	if request.Page <= 0 || request.PageSize <= 0 {
		result, err := c.categoryRepo.SearchWithFilters(c.db, filters, "created_at DESC")
		if err != nil {
			c.logger.Error("Failed to get categories", slog.String("error", err.Error()))
			return nil, err
		}

		response := &CategoriesResult{
			Categories: result,
			Total:      int64(len(result)),
			Page:       1,
			PageSize:   len(result),
			TotalPages: 1,
		}
		return response, nil
	}

	// ใช้ pagination พร้อม filters
	result, total, err := c.categoryRepo.SearchWithFiltersAndPagination(c.db, filters, "created_at DESC", request.Page, request.PageSize)
	if err != nil {
		c.logger.Error("Failed to get categories with pagination", slog.String("error", err.Error()))
		return nil, err
	}

	totalPages := int(total) / request.PageSize
	if int(total)%request.PageSize > 0 {
		totalPages++
	}

	response := &CategoriesResult{
		Categories: result,
		Total:      total,
		Page:       request.Page,
		PageSize:   request.PageSize,
		TotalPages: totalPages,
	}

	return response, nil
}
