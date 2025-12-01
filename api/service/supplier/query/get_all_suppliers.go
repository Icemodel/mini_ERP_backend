package query

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/model"

	"gorm.io/gorm"
)

type GetAllSuppliers struct {
	logger       *slog.Logger
	db           *gorm.DB
	SupplierRepo repository.SupplierRepository
}

type GetAllSuppliersRequest struct {
	OrderBy string `json:"order_by"`
}

type GetAllSuppliersResult struct {
	Suppliers []*model.Supplier `json:"suppliers"`
}

func NewGetAllSuppliersHandler(logger *slog.Logger, db *gorm.DB, repo repository.SupplierRepository) *GetAllSuppliers {
	return &GetAllSuppliers{
		logger:       logger,
		db:           db,
		SupplierRepo: repo,
	}
}

func (h *GetAllSuppliers) Handle(ctx context.Context, req *GetAllSuppliersRequest) (interface{}, error) {
	// Set default order by
	orderBy := req.OrderBy
	if orderBy == "" {
		orderBy = "created_at DESC"
	}

	// Get all suppliers from database
	suppliers, err := h.SupplierRepo.Searches(h.db, map[string]interface{}{}, orderBy)
	if err != nil {
		h.logger.Error("Failed to get all suppliers", "error", err)
		return nil, err
	}

	return GetAllSuppliersResult{Suppliers: suppliers}, nil
}
