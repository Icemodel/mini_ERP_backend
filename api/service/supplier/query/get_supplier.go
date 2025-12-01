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
	logger *slog.Logger
	db     *gorm.DB
	repo   repository.SupplierRepository
}

type GetSupplierRequest struct {
	SupplierId uuid.UUID `json:"supplier_id" validate:"required"`
}

func NewGetSupplierHandler(logger *slog.Logger, db *gorm.DB, repo repository.SupplierRepository) *GetSupplier {
	return &GetSupplier{
		logger: logger,
		db:     db,
		repo:   repo,
	}
}

func (h *GetSupplier) Handle(ctx context.Context, req GetSupplierRequest) (*model.Supplier, error) {
	// Find supplier by ID
	supplier, err := h.repo.Search(h.db, map[string]interface{}{
		"supplier_id": req.SupplierId,
	}, "")
	
	if err != nil {
		h.logger.Error("Failed to get supplier", "supplier_id", req.SupplierId, "error", err)
		return nil, err
	}

	return supplier, nil
}
