package query

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Stocks struct {
	logger               *slog.Logger
	db                   *gorm.DB
	stockTransactionRepo repository.StockTransaction
}

type StocksRequest struct {
	Page      int        `json:"page"`
	PageSize  int        `json:"page_size"`
	Search    string     `json:"search"`
	ProductId *uuid.UUID `json:"product_id"`
	SortBy    string     `json:"sort_by"`    // ฟิลด์ที่ต้องการ sort
	SortOrder string     `json:"sort_order"` // asc หรือ desc
}

type StocksResult struct {
	Stocks     []model.StockTransaction `json:"stocks"`
	Total      int64                    `json:"total"`
	Page       int                      `json:"page"`
	PageSize   int                      `json:"page_size"`
	TotalPages int                      `json:"total_pages"`
}

func NewStocks(logger *slog.Logger, db *gorm.DB, stockTransactionRepo repository.StockTransaction) *Stocks {
	return &Stocks{
		logger:               logger,
		db:                   db,
		stockTransactionRepo: stockTransactionRepo,
	}
}

func (s *Stocks) Handle(ctx context.Context, request StocksRequest) (*StocksResult, error) {
	filters := repository.StockTransactionSearchFilters{
		Search:    request.Search,
		ProductId: request.ProductId,
	}

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
		result, err := s.stockTransactionRepo.SearchWithFilters(s.db, filters, "created_at DESC")
		if err != nil {
			s.logger.Error("Failed to get products", slog.String("error", err.Error()))
			return nil, err
		}

		response := &StocksResult{
			Stocks:     result,
			Total:      int64(len(result)),
			Page:       1,
			PageSize:   len(result),
			TotalPages: 1,
		}
		return response, nil
	}

	stocks, total, err := s.stockTransactionRepo.SearchWithFiltersAndPagination(s.db, filters, orderBy, request.Page, request.PageSize)
	if err != nil {
		s.logger.Error("Failed to search stocks with filters and pagination", slog.String("error", err.Error()))
		return nil, err
	}

	totalPages := int((total + int64(request.PageSize) - 1) / int64(request.PageSize))

	result := &StocksResult{
		Stocks:     stocks,
		Total:      total,
		Page:       request.Page,
		PageSize:   request.PageSize,
		TotalPages: totalPages,
	}

	return result, nil
}