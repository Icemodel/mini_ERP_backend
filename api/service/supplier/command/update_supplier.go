package command

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateSupplier struct {
	logger       *slog.Logger
	db           *gorm.DB
	SupplierRepo repository.Supplier
}

type UpdateSupplierRequest struct {
	SupplierId uuid.UUID `json:"-"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
	Address    string    `json:"address"`
}

type UpdateSupplierResult struct {
	Supplier model.Supplier `json:"supplier"`
}

func NewUpdateSupplier(logger *slog.Logger, db *gorm.DB, repo repository.Supplier) *UpdateSupplier {
	return &UpdateSupplier{
		logger:       logger,
		db:           db,
		SupplierRepo: repo,
	}
}

func (h *UpdateSupplier) Handle(ctx context.Context, cmd *UpdateSupplierRequest) (*UpdateSupplierResult, error) {

	// Validate input
	if cmd.SupplierId == uuid.Nil {
		h.logger.Error("Supplier ID is required")
		return nil, gorm.ErrInvalidData
	}

	// Check if at least one field is provided for update
	if cmd.Name == "" && cmd.Phone == "" && cmd.Email == "" && cmd.Address == "" {
		h.logger.Warn("No fields to update")
		return nil, gorm.ErrInvalidData
	}

	// Check if supplier exists 
	supplier_id := map[string]interface{}{
		"supplier_id": cmd.SupplierId,
	}
	existingSupplier, err := h.SupplierRepo.Search(
		h.db.WithContext(ctx),
		supplier_id,
		"",
	)
	if err != nil {
		h.logger.Error("Supplier not found", "supplier_id", cmd.SupplierId)
		return nil, err
	}

	if cmd.Name != "" {
		existingSupplier.Name = cmd.Name
	}

	if cmd.Phone != "" {
		existingSupplier.Phone = cmd.Phone
	}
	
	if cmd.Email != "" {
		existingSupplier.Email = cmd.Email
	}

	if cmd.Address != "" {
		existingSupplier.Address = cmd.Address
	}

	// Begin transaction 
	tx := h.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

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

	return &UpdateSupplierResult{
		Supplier: *existingSupplier,
	}, nil
}
