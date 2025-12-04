package command

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeleteById struct {
	logger      *slog.Logger
	db          *gorm.DB
	productRepo repository.Product
}

type DeleteByIdRequest struct {
	ProductId uuid.UUID `json:"product_id"`
}

type DeleteByIdResult struct {
	Deleted bool   `json:"deleted"`
	Message string `json:"message,omitempty"`
}

func NewDeleteById(logger *slog.Logger, db *gorm.DB, productRepo repository.Product) *DeleteById {
	return &DeleteById{
		logger:      logger,
		db:          db,
		productRepo: productRepo,
	}
}

func (d *DeleteById) Handle(ctx context.Context, request DeleteByIdRequest) (*DeleteByIdResult, error) {
	err := d.productRepo.DeleteById(d.db, request.ProductId)
	if err != nil {
		d.logger.Error("Failed to delete product by id", slog.String("error", err.Error()))
		return nil, err
	}
	response := &DeleteByIdResult{
		Deleted: true,
		Message: "Product deleted successfully",
	}
	return response, nil
}
