package command

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateSupplier struct {
	logger *slog.Logger
	db     *gorm.DB
	repo   repository.SupplierRepository
}

type UpdateSupplierRequest struct {
	SupplierId uuid.UUID `json:"supplier_id" validate:"required"`
	Name       string    `json:"name" validate:"required"`
	Phone      string    `json:"phone" validate:"required"`
	Email      string    `json:"email" validate:"required,email"`
	Address    string    `json:"address" validate:"required"`
}

func NewUpdateSupplierHandler(logger *slog.Logger, db *gorm.DB, repo repository.SupplierRepository) *UpdateSupplier {
	return &UpdateSupplier{
		logger: logger,
		db:     db,
		repo:   repo,
	}
}

func (h *UpdateSupplier) Handle(ctx context.Context, cmd UpdateSupplierRequest) error {
	// Check if supplier exists
	existingSupplier, err := h.repo.Search(h.db, map[string]interface{}{
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

	// Update supplier fields
	existingSupplier.Name = cmd.Name
	existingSupplier.Phone = cmd.Phone
	existingSupplier.Email = cmd.Email
	existingSupplier.Address = cmd.Address

	// Save to database
	if err := h.repo.UpdateBySupplierId(tx, cmd.SupplierId, existingSupplier); err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		h.logger.Error("Failed to commit transaction", "error", err)
		return err
	}

	h.logger.Info("Supplier updated successfully", "supplier_id", cmd.SupplierId)
	return nil
}
