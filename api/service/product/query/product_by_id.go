package query

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductById struct {
	logger      *slog.Logger
	db          *gorm.DB
	productRepo repository.Product
}

type ProductByIdRequest struct {
	ProductId uuid.UUID `json:"product_id" binding:"required"`
}

type ProductResult struct {
	model.Product `json:"product"`
}

func NewProductById(logger *slog.Logger, db *gorm.DB, productRepo repository.Product) *ProductById {
	return &ProductById{
		logger:      logger,
		db:          db,
		productRepo: productRepo,
	}
}

func (p *ProductById) Handle(ctx context.Context, request ProductByIdRequest) (*ProductResult, error) {
	conditions := map[string]interface{}{
		"product_id": request.ProductId,
	}

	result, err := p.productRepo.Search(p.db, conditions, "")

	if err != nil {
		p.logger.Error("Failed to get product by id", slog.String("error", err.Error()))
		return nil, err
	}

	response := &ProductResult{
		Product: *result,
	}

	return response, nil
}
