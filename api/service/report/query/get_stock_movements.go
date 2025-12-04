package query

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"
	"time"

	"gorm.io/gorm"
)

type StockMovements struct {
	logger     *slog.Logger
	db         *gorm.DB
	reportRepo repository.Report
}

type StockMovementsRequest struct {
	FromDate time.Time `json:"from_date"`
	ToDate   time.Time `json:"to_date"`
}

type StockMovementsResult struct {
	Movements []repository.StockMovementResult `json:"movements"`
}

func NewStockMovements(
	logger *slog.Logger,
	db *gorm.DB,
	reportRepo repository.Report,
) *StockMovements {
	return &StockMovements{
		logger:     logger,
		db:         db,
		reportRepo: reportRepo,
	}
}

func (h *StockMovements) Handle(ctx context.Context, req *StockMovementsRequest) (*StockMovementsResult, error) {
	movements, err := h.reportRepo.GetStockMovements(h.db, req.FromDate, req.ToDate)
	if err != nil {
		h.logger.Error("Failed to get stock movements", "error", err)
		return nil, err
	}

	return &StockMovementsResult{Movements: movements}, nil
}
