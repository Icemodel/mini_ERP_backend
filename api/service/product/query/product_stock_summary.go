package query

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductStockSummary struct {
	logger               *slog.Logger
	db                   *gorm.DB
	productRepo          repository.Product
	stockTransactionRepo repository.StockTransaction
}

type StockSummary struct {
	TotalIn      int64 `json:"total_in"`
	TotalOut     int64 `json:"total_out"`
	TotalAdjust  int64 `json:"total_adjust"`
	CurrentStock int64 `json:"current_stock"`
}

type ProductStockSummaryRequest struct {
	ProductId uuid.UUID `json:"product_id"`
}

type ProductStockSummaryResult struct {
	Product    model.Product `json:"product"`
	Stock      StockSummary  `json:"stock_summary"`
	MinStock   int64         `json:"min_stock"`
	IsLowStock bool          `json:"is_low_stock"`
}

func NewProductStockSummary(logger *slog.Logger, db *gorm.DB, productRepo repository.Product, stockTransactionRepo repository.StockTransaction) *ProductStockSummary {
	return &ProductStockSummary{
		logger:               logger,
		db:                   db,
		productRepo:          productRepo,
		stockTransactionRepo: stockTransactionRepo,
	}
}

func (p *ProductStockSummary) Handle(ctx context.Context, request ProductStockSummaryRequest) (*ProductStockSummaryResult, error) {
	conditions := map[string]interface{}{
		"product_id": request.ProductId,
	}

	product, err := p.productRepo.Search(p.db, conditions, "")
	if err != nil {
		p.logger.Error("Failed to get product by id", slog.String("error", err.Error()))
		return nil, err
	}

	totalIn, totalOut, totalAdjust, err := p.stockTransactionRepo.StockSummary(p.db, request.ProductId)
	if err != nil {
		p.logger.Error("Failed to get stock summary", slog.String("error", err.Error()))
		return nil, err
	}

	currentStock := totalIn - totalOut + totalAdjust

	response := &ProductStockSummaryResult{
		Product: *product,
		Stock: StockSummary{
			TotalIn:      totalIn,
			TotalOut:     totalOut,
			TotalAdjust:  totalAdjust,
			CurrentStock: currentStock,
		},
		MinStock:   product.MinStock,
		IsLowStock: currentStock < product.MinStock,
	}

	return response, nil
}