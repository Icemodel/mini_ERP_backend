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
	logger       *slog.Logger
	db           *gorm.DB
	categoryRepo repository.Category
}

type UpdateRequest struct {
	CategoryId  uuid.UUID `json:"category_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
}

type UpdateResult struct {
	Category model.Category `json:"category"`
}

func NewUpdate(logger *slog.Logger, db *gorm.DB, categoryRepo repository.Category) *Update {
	return &Update{
		logger:       logger,
		db:           db,
		categoryRepo: categoryRepo,
	}
}

func (u *Update) Handle(ctx context.Context, request UpdateRequest) (*UpdateResult, error) {
	condition := map[string]interface{}{
		"category_id": request.CategoryId,
	}

	// ตรวจสอบว่า category มีอยู่หรือไม่
	category, err := u.categoryRepo.Search(u.db, condition, "")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u.logger.Error("Category not found", slog.String("category_id", request.CategoryId.String()))
			return nil, errors.New("category not found")
		}
		u.logger.Error("Failed to get category", slog.String("error", err.Error()))
		return nil, err
	}

	// ตรวจสอบชื่อซ้ำ (ยกเว้นตัวเอง)
	existed, err := u.categoryRepo.ExitedByNameExcludeId(u.db, request.Name, request.CategoryId)
	if err != nil {
		u.logger.Error("Failed to check if category name exists", slog.String("error", err.Error()))
		return nil, err
	}

	if existed {
		u.logger.Error("Category name already exists", slog.String("name", request.Name))
		return nil, errors.New("category name already exists")
	}

	// อัพเดทข้อมูล
	category.Name = request.Name
	category.Description = request.Description
	category.UpdatedAt = time.Now()

	if err := u.categoryRepo.Update(u.db, category); err != nil {
		u.logger.Error("Failed to update category", slog.String("error", err.Error()))
		return nil, err
	}

	response := &UpdateResult{
		Category: *category,
	}

	return response, nil
}