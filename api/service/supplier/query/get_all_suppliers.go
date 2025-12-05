package query

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/model"

	"gorm.io/gorm"
)

type AllSuppliers struct {
	logger       *slog.Logger
	db           *gorm.DB
	SupplierRepo repository.Supplier
}

type AllSuppliersRequest struct {
	SortOrder string `json:"sort_order"`
}

type AllSuppliersResult struct {
	Suppliers []*model.Supplier `json:"suppliers"`
}

func NewAllSuppliers(logger *slog.Logger, db *gorm.DB, repo repository.Supplier) *AllSuppliers {
	return &AllSuppliers{
		logger:       logger,
		db:           db,
		SupplierRepo: repo,
	}
}

func (h *AllSuppliers) Handle(ctx context.Context, req *AllSuppliersRequest) (*AllSuppliersResult, error) {
	// Build orderBy string
	orderBy := "created_at DESC" // default
	if req.SortOrder == "asc" || req.SortOrder == "ASC" {
		orderBy = "created_at ASC"
	}

	// Get all suppliers from database (ใช้ context)
	suppliers, err := h.SupplierRepo.Searches(
		h.db.WithContext(ctx),
		map[string]interface{}{},
		orderBy,
	)
	if err != nil {
		h.logger.Error("Failed to get all suppliers", "error", err)
		return nil, err
	}

	return &AllSuppliersResult{Suppliers: suppliers}, nil
}
