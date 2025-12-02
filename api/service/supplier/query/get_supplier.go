package query

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Supplier struct {
	logger       *slog.Logger
	db           *gorm.DB
	SupplierRepo repository.Supplier
}

type SupplierRequest struct {
	SupplierId uuid.UUID `json:"supplier_id" validate:"required"`
}

type SupplierResult struct {
	Supplier *model.Supplier `json:"supplier"`
}

func NewSupplier(logger *slog.Logger, db *gorm.DB, repo repository.Supplier) *Supplier {
	return &Supplier{
		logger:       logger,
		db:           db,
		SupplierRepo: repo,
	}
}

func (h *Supplier) Handle(ctx context.Context, req *SupplierRequest) (*SupplierResult, error) {
	// Find supplier by ID
	supplier_id := map[string]interface{}{
		"supplier_id": req.SupplierId,
	}
	supplier, err := h.SupplierRepo.Search(h.db, supplier_id, "")

	if err != nil {
		h.logger.Error("Failed to get supplier", "supplier_id", req.SupplierId, "error", err)
		return nil, err
	}

	return &SupplierResult{Supplier: supplier}, nil
}
