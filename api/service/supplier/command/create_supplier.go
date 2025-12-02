package command

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type CreateSupplier struct {
	logger *slog.Logger
	db     *gorm.DB
	SupplierRepo   repository.Supplier
}

type CreateSupplierRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Address string `json:"address"`
}


func NewCreateSupplier(logger *slog.Logger, db *gorm.DB, repo repository.Supplier) *CreateSupplier {
	return &CreateSupplier{
		logger:       logger,
		db:           db,
		SupplierRepo: repo,
	}
}

func (h *CreateSupplier) Handle(ctx context.Context, cmd *CreateSupplierRequest) (interface{}, error) {

	// Check if email already exists
	email := map[string]interface{}{
		"email": cmd.Email,
	}

	existingSupplier, err := h.SupplierRepo.Search(h.db, email, "")
	
	if err == nil && existingSupplier != nil {
		h.logger.Warn("Supplier with this email already exists", "email", cmd.Email)
		return nil, gorm.ErrDuplicatedKey
	}

	// Begin transaction
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create supplier model
	supplier := &model.Supplier{
		SupplierId: uuid.New(),
		Name:       cmd.Name,
		Phone:      cmd.Phone,
		Email:      cmd.Email,
		Address:    cmd.Address,
	}

	// Save to database
	if err := h.SupplierRepo.Create(tx, supplier); err != nil {
		tx.Rollback()
		h.logger.Error("Failed to create supplier", "error", err)
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		h.logger.Error("Failed to commit transaction", "error", err)
		return nil, err
	}

	h.logger.Info("Supplier created successfully", "supplier_id", supplier.SupplierId)
	return supplier, nil
}
