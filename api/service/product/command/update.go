package command

import (
	"context"
	"errors"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Update struct {
	logger      *slog.Logger
	db          *gorm.DB
	productRepo repository.Product
}

type UpdateRequest struct {
	ProductId    uuid.UUID `json:"product_id"`
	CategoryId   uuid.UUID `json:"category_id"`
	Name         string    `json:"name"`
	CostPrice    float64   `json:"cost_price"`
	SellingPrice float64   `json:"selling_price"`
	Unit         string    `json:"unit"`
	MinStock     int64     `json:"min_stock"`
}

type UpdateResult struct {
	Product model.Product `json:"product"`
}

func NewUpdate(logger *slog.Logger, db *gorm.DB, productRepo repository.Product) *Update {
	return &Update{
		logger:      logger,
		db:          db,
		productRepo: productRepo,
	}
}

func (u *Update) Handle(ctx context.Context, request UpdateRequest) (*UpdateResult, error) {
	condition := map[string]interface{}{
		"product_id": request.ProductId,
	}

	product, err := u.productRepo.Search(u.db, condition, "")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u.logger.Error("Product not found", slog.String("product_id", request.ProductId.String()))
			return nil, errors.New("product not found")
		}
		u.logger.Error("Failed to get product", slog.String("error", err.Error()))
		return nil, err
	}

	product.Name = request.Name
	product.CategoryId = request.CategoryId
	product.CostPrice = request.CostPrice
	product.SellingPrice = request.SellingPrice
	product.Unit = request.Unit
	product.MinStock = request.MinStock
	product.UpdatedAt = time.Now()

	if err := u.productRepo.Update(u.db, product); err != nil {
		u.logger.Error("Failed to update product", slog.String("error", err.Error()))
		return nil, err
	}

	response := &UpdateResult{
		Product: *product,
	}

	return response, nil
}
