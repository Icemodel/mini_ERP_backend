package command

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeleteSupplier struct {
	logger       *slog.Logger
	db           *gorm.DB
	SupplierRepo repository.Supplier
}

type DeleteSupplierRequest struct {
	SupplierId uuid.UUID `json:"supplier_id"`
}

type DeleteSupplierResult struct {
	Deleted bool `json:"deleted"`
	Message string `json:"message,omitempty"`
}

func NewDeleteSupplier(logger *slog.Logger, db *gorm.DB, repo repository.Supplier) *DeleteSupplier {
	return &DeleteSupplier{
		logger:       logger,
		db:           db,
		SupplierRepo: repo,
	}
}

func (h *DeleteSupplier) Handle(ctx context.Context, cmd *DeleteSupplierRequest) (*DeleteSupplierResult, error) {
	// Validate input
	if cmd.SupplierId == uuid.Nil {
		h.logger.Error("Supplier ID is required")
		return nil, gorm.ErrInvalidData
	}

	// Check if supplier exists 
	supplier_id := map[string]interface{}{
		"supplier_id": cmd.SupplierId,
	}
	supplier, err := h.SupplierRepo.Search(
		h.db.WithContext(ctx),
		supplier_id,
		"",
	)
	if err != nil {
		h.logger.Error("Supplier not found", "supplier_id", cmd.SupplierId)
		return nil, err
	}

	// Delete from database 
	if err := h.SupplierRepo.Delete(h.db.WithContext(ctx), supplier); err != nil {
		h.logger.Error("Failed to delete supplier", "error", err)
		return nil, err
	}

	return &DeleteSupplierResult{
		Deleted: true,
		Message: "Supplier deleted successfully",
	}, nil
}
