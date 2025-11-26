package query

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryById struct {
	logger       *slog.Logger
	db           *gorm.DB
	categoryRepo repository.Category
}

type CategoryByIdRequest struct {
	CategoryId uuid.UUID `json:"category_id"`
}

type CategoryByIdResult struct {
	model.Category
}

func NewCategoryById(logger *slog.Logger, db *gorm.DB, categoryRepo repository.Category) *CategoryById {
	return &CategoryById{
		logger:       logger,
		db:           db,
		categoryRepo: categoryRepo,
	}
}

func (c *CategoryById) Handle(ctx context.Context, request CategoryByIdRequest) (*CategoryByIdResult, error) {
	conditions := map[string]interface{}{
		"category_id": request.CategoryId,
	}
	result, err := c.categoryRepo.Search(c.db, conditions, "")
	if err != nil {
		c.logger.Error("Failed to get category by id", slog.String("error", err.Error()))
		return nil, err
	}
	response := &CategoryByIdResult{
		Category: *result,
	}
	return response, nil
}
