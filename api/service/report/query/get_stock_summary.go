package query

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"

	"gorm.io/gorm"
)

type StockSummary struct {
	logger     *slog.Logger
	db         *gorm.DB
	reportRepo repository.Report
}

type StockSummaryRequest struct{}

type StockSummaryResult struct {
	Products          []repository.StockSummaryResult `json:"products"`
	TotalStockOnHand  int64                           `json:"total_stock_on_hand"`
	TotalCostValue    float64                         `json:"total_cost_value"`
	TotalSellingValue float64                         `json:"total_selling_value"`
	LowStock          []repository.StockSummaryResult `json:"low_stock"`
	LowStockCount     int                             `json:"low_stock_count"`
}

func NewStockSummary(
	logger *slog.Logger,
	db *gorm.DB,
	reportRepo repository.Report,
) *StockSummary {
	return &StockSummary{
		logger:     logger,
		db:         db,
		reportRepo: reportRepo,
	}
}

func (h *StockSummary) Handle(ctx context.Context, req *StockSummaryRequest) (*StockSummaryResult, error) {
	products, err := h.reportRepo.GetStockSummary(h.db)
	if err != nil {
		h.logger.Error("Failed to get stock summary", "error", err)
		return nil, err
	}

	var totalStock int64
	var totalCost float64
	var totalSelling float64
	var lowStock []repository.StockSummaryResult

	for i := range products {
		p := &products[i]
		totalStock += p.StockOnHand
		totalCost += p.TotalCostValue
		totalSelling += p.TotalSellingValue
		if p.StockOnHand < p.MinStock {
			p.IsLowStock = true
			lowStock = append(lowStock, *p)
		}
	}

	return &StockSummaryResult{
		Products:          products,
		TotalStockOnHand:  totalStock,
		TotalCostValue:    totalCost,
		TotalSellingValue: totalSelling,
		LowStock:          lowStock,
		LowStockCount:     len(lowStock),
	}, nil
}