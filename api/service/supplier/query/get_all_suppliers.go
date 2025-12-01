package query

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/model"

	"gorm.io/gorm"
)

type GetAllSuppliers struct {
	logger *slog.Logger
	db     *gorm.DB
	repo   repository.SupplierRepository
}

type GetAllSuppliersRequest struct {
	OrderBy string `json:"order_by"`
}

func NewGetAllSuppliersHandler(logger *slog.Logger, db *gorm.DB, repo repository.SupplierRepository) *GetAllSuppliers {
	return &GetAllSuppliers{
		logger: logger,
		db:     db,
		repo:   repo,
	}
}

func (h *GetAllSuppliers) Handle(ctx context.Context, req GetAllSuppliersRequest) ([]*model.Supplier, error) {
	// Set default order by
	orderBy := req.OrderBy
	if orderBy == "" {
		orderBy = "created_at DESC"
	}

	// Get all suppliers from database
	suppliers, err := h.repo.Searches(h.db, map[string]interface{}{}, orderBy)
	if err != nil {
		h.logger.Error("Failed to get all suppliers", "error", err)
		return nil, err
	}

	return suppliers, nil
}
