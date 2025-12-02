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
	OrderBy string `json:"order_by"`
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
	// Set default order by
	orderBy := req.OrderBy
	if orderBy == "" {
		orderBy = "created_at DESC"
	}

	// Get all suppliers from database
	conditions := map[string]interface{}{}
	suppliers, err := h.SupplierRepo.Searches(h.db, conditions, orderBy)
	if err != nil {
		h.logger.Error("Failed to get all suppliers", "error", err)
		return nil, err
	}

	return &AllSuppliersResult{Suppliers: suppliers}, nil
}
