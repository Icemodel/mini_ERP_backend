package query

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"
	"time"

	"gorm.io/gorm"
)

type GetStockMovements struct {
	logger     *slog.Logger
	db         *gorm.DB
	reportRepo repository.Report
}

type GetStockMovementsRequest struct {
	FromDate time.Time `json:"from_date"`
	ToDate   time.Time `json:"to_date"`
}

type GetStockMovementsResult struct {
	Movements []repository.StockMovementResult `json:"movements"`
}

func NewGetStockMovementsHandler(
	logger *slog.Logger,
	db *gorm.DB,
	reportRepo repository.Report,
) *GetStockMovements {
	return &GetStockMovements{
		logger:     logger,
		db:         db,
		reportRepo: reportRepo,
	}
}

func (h *GetStockMovements) Handle(ctx context.Context, req *GetStockMovementsRequest) (*GetStockMovementsResult, error) {
	movements, err := h.reportRepo.GetStockMovements(h.db, req.FromDate, req.ToDate)
	if err != nil {
		h.logger.Error("Failed to get stock movements", "error", err)
		return nil, err
	}

	return &GetStockMovementsResult{Movements: movements}, nil
}
