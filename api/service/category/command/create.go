package command

import (
	"context"
	"errors"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/model"
	"time"

	"gorm.io/gorm"
)

type Create struct {
	logger       *slog.Logger
	db           *gorm.DB
	categoryRepo repository.Category
}

type CreateRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type CreateResult struct {
	Category model.Category `json:"category"`
}

func NewCreate(logger *slog.Logger, db *gorm.DB, categoryRepo repository.Category) *Create {
	return &Create{
		logger:       logger,
		db:           db,
		categoryRepo: categoryRepo,
	}
}

func (c *Create) Handle(ctx context.Context, request CreateRequest) (*CreateResult, error) {
	// ตรวจสอบชื่อซ้ำ
	existed, err := c.categoryRepo.ExitedByName(c.db, request.Name)
	if err != nil {
		c.logger.Error("Failed to check if category name exists", slog.String("error", err.Error()))
		return nil, err
	}

	if existed {
		c.logger.Error("Category name already exists", slog.String("name", request.Name))
		return nil, errors.New("category name already exists")
	}

	// สร้าง category ใหม่
	category := &model.Category{
		// CategoryId:  uuid.New(),
		Name:        request.Name,
		Description: request.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := c.categoryRepo.Create(c.db, category); err != nil {
		c.logger.Error("Failed to create category", slog.String("error", err.Error()))
		return nil, err
	}

	response := &CreateResult{
		Category: *category,
	}

	return response, nil
}