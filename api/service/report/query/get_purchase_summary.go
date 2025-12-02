package query

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"

	"gorm.io/gorm"
)

type PurchaseSummary struct {
	logger     *slog.Logger
	db         *gorm.DB
	reportRepo repository.Report
}

type PurchaseSummaryRequest struct {
	Year  int `json:"year"`
	Month int `json:"month"`
}

type PurchaseSummaryResult struct {
	Summary        []repository.PurchaseSummaryResult `json:"summary"`
	TotalOrders    int64                              `json:"total_orders"`
	TotalAmount    uint64                             `json:"total_amount"`
	ReceivedOrders int64                              `json:"received_orders"`
	ReceivedAmount uint64                             `json:"received_amount"`
}

func NewPurchaseSummary(
	logger *slog.Logger,
	db *gorm.DB,
	reportRepo repository.Report,
) *PurchaseSummary {
	return &PurchaseSummary{
		logger:     logger,
		db:         db,
		reportRepo: reportRepo,
	}
}

func (h *PurchaseSummary) Handle(ctx context.Context, req *PurchaseSummaryRequest) (*PurchaseSummaryResult, error) {
	summary, err := h.reportRepo.GetPurchaseSummary(h.db, req.Year, req.Month)
	if err != nil {
		h.logger.Error("Failed to get purchase summary", "error", err)
		return nil, err
	}

	var totalOrders int64
	var totalAmount uint64
	var receivedOrders int64
	var receivedAmount uint64

	for _, s := range summary {
		totalOrders += s.TotalOrders
		totalAmount += s.TotalAmount
		if s.Status == "RECEIVED" {
			receivedOrders = s.TotalOrders
			receivedAmount = s.TotalAmount
		}
	}

	return &PurchaseSummaryResult{
		Summary:        summary,
		TotalOrders:    totalOrders,
		TotalAmount:    totalAmount,
		ReceivedOrders: receivedOrders,
		ReceivedAmount: receivedAmount,
	}, nil
}
