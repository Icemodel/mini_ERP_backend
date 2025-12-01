package command

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeleteSupplier struct {
	logger *slog.Logger
	db     *gorm.DB
	repo   repository.SupplierRepository
}

type DeleteSupplierRequest struct {
	SupplierId uuid.UUID `json:"supplier_id" validate:"required"`
}

func NewDeleteSupplierHandler(logger *slog.Logger, db *gorm.DB, repo repository.SupplierRepository) *DeleteSupplier {
	return &DeleteSupplier{
		logger: logger,
		db:     db,
		repo:   repo,
	}
}

func (h *DeleteSupplier) Handle(ctx context.Context, cmd DeleteSupplierRequest) error {
	// Check if supplier exists
	supplier, err := h.repo.Search(h.db, map[string]interface{}{
		"supplier_id": cmd.SupplierId,
	}, "")
	if err != nil {
		h.logger.Error("Supplier not found", "supplier_id", cmd.SupplierId)
		return err
	}

	// Begin transaction
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete from database
	if err := h.repo.Delete(tx, supplier); err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		h.logger.Error("Failed to commit transaction", "error", err)
		return err
	}

	h.logger.Info("Supplier deleted successfully", "supplier_id", cmd.SupplierId)
	return nil
}
