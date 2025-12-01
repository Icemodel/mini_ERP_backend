package command

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateSupplier struct {
	logger       *slog.Logger
	db           *gorm.DB
	SupplierRepo repository.Supplier
}

type UpdateSupplierRequest struct {
	SupplierId uuid.UUID `json:"supplier_id" validate:"required"`
	Name       string    `json:"name" validate:"required"`
	Phone      string    `json:"phone" validate:"required"`
	Email      string    `json:"email" validate:"required,email"`
	Address    string    `json:"address" validate:"required"`
}

func NewUpdateSupplierHandler(logger *slog.Logger, db *gorm.DB, repo repository.Supplier) *UpdateSupplier {
	return &UpdateSupplier{
		logger:       logger,
		db:           db,
		SupplierRepo: repo,
	}
}

func (h *UpdateSupplier) Handle(ctx context.Context, cmd *UpdateSupplierRequest) (interface{}, error) {
	// Check if supplier exists
	supplier_id := map[string]interface{}{
		"supplier_id": cmd.SupplierId,
	}
	existingSupplier, err := h.SupplierRepo.Search(h.db, supplier_id, "")
	if err != nil {
		h.logger.Error("Supplier not found", "supplier_id", cmd.SupplierId)
		return nil, err
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
	if err := h.SupplierRepo.UpdateBySupplierId(tx, cmd.SupplierId, existingSupplier); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		h.logger.Error("Failed to commit transaction", "error", err)
		return nil, err
	}

	h.logger.Info("Supplier updated successfully", "supplier_id", cmd.SupplierId)
	return nil, nil
}
