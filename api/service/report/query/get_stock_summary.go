package query

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"

	"gorm.io/gorm"
)

type GetStockSummary struct {
	logger     *slog.Logger
	db         *gorm.DB
	reportRepo repository.ReportRepository
}

type GetStockSummaryRequest struct{}

type GetStockSummaryResult struct {
	Products []repository.StockSummaryResult `json:"products"`
}

func NewGetStockSummaryHandler(
	logger *slog.Logger,
	db *gorm.DB,
	reportRepo repository.ReportRepository,
) *GetStockSummary {
	return &GetStockSummary{
		logger:     logger,
		db:         db,
		reportRepo: reportRepo,
	}
}

func (h *GetStockSummary) Handle(ctx context.Context, req *GetStockSummaryRequest) (*GetStockSummaryResult, error) {
	products, err := h.reportRepo.GetStockSummary(h.db)
	if err != nil {
		h.logger.Error("Failed to get stock summary", "error", err)
		return nil, err
	}

	return &GetStockSummaryResult{Products: products}, nil
}
