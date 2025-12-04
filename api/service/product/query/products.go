package query

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Products struct {
	logger      *slog.Logger
	db          *gorm.DB
	productRepo repository.Product
}

type ProductsRequest struct {
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
	Search     string     `json:"search"`
	CategoryId *uuid.UUID `json:"category_id"`
	SortBy     string     `json:"sort_by"`    // ฟิลด์ที่ต้องการ sort
	SortOrder  string     `json:"sort_order"` // asc หรือ desc
}

type ProductsResult struct {
	Products   []model.Product `json:"products"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

func NewProducts(logger *slog.Logger, db *gorm.DB, productRepo repository.Product) *Products {
	return &Products{
		logger:      logger,
		db:          db,
		productRepo: productRepo,
	}
}

func (p *Products) Handle(ctx context.Context, request ProductsRequest) (*ProductsResult, error) {
	filters := repository.ProductSearchFilters{
		Search:     request.Search,
		CategoryId: request.CategoryId,
	}

	// สร้าง orderBy string
	orderBy := "created_at DESC" // default
	if request.SortBy != "" {
		// กำหนดฟิลด์ที่อนุญาตให้ sort
		allowedSortFields := map[string]bool{
			"name":          true,
			"cost_price":    true,
			"selling_price": true,
			"unit":          true,
			"min_stock":     true,
			"created_at":    true,
			"updated_at":    true,
		}

		if allowedSortFields[request.SortBy] {
			sortOrder := "DESC"
			if request.SortOrder == "asc" || request.SortOrder == "ASC" {
				sortOrder = "ASC"
			}
			orderBy = request.SortBy + " " + sortOrder
		}
	}

	if request.Page <= 0 || request.PageSize <= 0 {
		result, err := p.productRepo.SearchWithFilters(p.db, filters, "created_at DESC")
		if err != nil {
			p.logger.Error("Failed to get products", slog.String("error", err.Error()))
			return nil, err
		}

		response := &ProductsResult{
			Products:   result,
			Total:      int64(len(result)),
			Page:       1,
			PageSize:   len(result),
			TotalPages: 1,
		}
		return response, nil
	}

	// ดึงข้อมูลแบบ pagination เสมอ
	result, total, err := p.productRepo.SearchWithFiltersAndPagination(p.db, filters, orderBy, request.Page, request.PageSize)
	if err != nil {
		p.logger.Error("Failed to get products", slog.String("error", err.Error()))
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

	response := &ProductsResult{
		Products:   result,
		Total:      total,
		Page:       request.Page,
		PageSize:   request.PageSize,
		TotalPages: totalPages,
	}

	return response, nil
}