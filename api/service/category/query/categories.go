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
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
	Search    string `json:"search"`     // ค้นหาทั้งชื่อและคำอธิบาย
	SortBy    string `json:"sort_by"`    // ฟิลด์ที่ต้องการ sort
	SortOrder string `json:"sort_order"` // asc หรือ desc
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

	// สร้าง orderBy string
	orderBy := "created_at DESC" // default
	if request.SortBy != "" {
		// กำหนดฟิลด์ที่อนุญาตให้ sort
		allowedSortFields := map[string]bool{
			"name":       true,
			"created_at": true,
			"updated_at": true,
		}

		if allowedSortFields[request.SortBy] {
			sortOrder := "DESC"
			if request.SortOrder == "asc" || request.SortOrder == "ASC" {
				sortOrder = "ASC"
			}
			orderBy = request.SortBy + " " + sortOrder
		}
	}

	// ดึงข้อมูลแบบ pagination เสมอ
	result, total, err := c.categoryRepo.SearchWithFiltersAndPagination(c.db, filters, orderBy, request.Page, request.PageSize)
	if err != nil {
		c.logger.Error("Failed to get categories with pagination", slog.String("error", err.Error()))
		return nil, err
	}

	// คำนวณจำนวนหน้าทั้งหมด
	totalPages := 0
	if total > 0 {
		totalPages = int(total) / request.PageSize
		if int(total)%request.PageSize > 0 {
			totalPages++
		}
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
