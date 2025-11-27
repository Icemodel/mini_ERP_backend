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

type Create struct {
	logger      *slog.Logger
	db          *gorm.DB
	productRepo repository.Product
}

type CreateRequest struct {
	ProductCode  string    `json:"product_code"`
	CategoryId   uuid.UUID `json:"category_id"`
	Name         string    `json:"name"`
	CostPrice    float64   `json:"cost_price"`
	SellingPrice float64   `json:"selling_price"`
	Unit         int64     `json:"unit"`
	MinStock     int64     `json:"min_stock"`
}

type CreateResult struct {
	Product model.Product `json:"product"`
}

func NewCreate(logger *slog.Logger, db *gorm.DB, productRepo repository.Product) *Create {
	return &Create{
		logger:      logger,
		db:          db,
		productRepo: productRepo,
	}
}

func (c *Create) Handle(ctx context.Context, request CreateRequest) (*CreateResult, error) {
	existed, err := c.productRepo.ExitedByProductCode(c.db, request.ProductCode)

	if err != nil {
		c.logger.Error("Failed to check if product code exists", slog.String("error", err.Error()))
		return nil, err
	}

	if existed {
		c.logger.Error("Product code already exists", slog.String("product_code", request.ProductCode))
		return nil, errors.New("product code already exists")
	}

	product := &model.Product{
		ProductId:    uuid.New(),
		ProductCode:  request.ProductCode,
		CategoryId:   request.CategoryId,
		Name:         request.Name,
		CostPrice:    request.CostPrice,
		SellingPrice: request.SellingPrice,
		Unit:         request.Unit,
		MinStock:     request.MinStock,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Category:     nil, // ไม่ load category
	}

	if err := c.productRepo.Create(c.db, product); err != nil {
		c.logger.Error("Failed to create product", slog.String("error", err.Error()))
		return nil, err
	}

	response := &CreateResult{
		Product: *product,
	}

	return response, nil
}
