package query

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetSupplier struct {
	logger       *slog.Logger
	db           *gorm.DB
	SupplierRepo repository.SupplierRepository
}

type GetSupplierRequest struct {
	SupplierId uuid.UUID `json:"supplier_id" validate:"required"`
}

type GetSupplierResult struct {
	Supplier *model.Supplier `json:"supplier"`
}

func NewGetSupplierHandler(logger *slog.Logger, db *gorm.DB, repo repository.SupplierRepository) *GetSupplier {
	return &GetSupplier{
		logger:       logger,
		db:           db,
		SupplierRepo: repo,
	}
}

func (h *GetSupplier) Handle(ctx context.Context, req *GetSupplierRequest) (interface{}, error) {
	// Find supplier by ID
	supplier, err := h.SupplierRepo.Search(h.db, map[string]interface{}{
		"supplier_id": req.SupplierId,
	}, "")
	
	if err != nil {
		h.logger.Error("Failed to get supplier", "supplier_id", req.SupplierId, "error", err)
		return nil, err
	}

	return GetSupplierResult{Supplier: supplier}, nil
}
