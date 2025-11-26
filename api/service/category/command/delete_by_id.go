package command

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeleteById struct {
	logger       *slog.Logger
	db           *gorm.DB
	categoryRepo repository.Category
}

type DeleteByIdRequest struct {
	CategoryId uuid.UUID `json:"category_id"`
}

type DeleteByIdResult struct {
	Deleted bool   `json:"deleted"`
	Message string `json:"message,omitempty"`
}

func NewDeleteById(logger *slog.Logger, db *gorm.DB, categoryRepo repository.Category) *DeleteById {
	return &DeleteById{
		logger:       logger,
		db:           db,
		categoryRepo: categoryRepo,
	}
}

func (d *DeleteById) Handle(ctx context.Context, request DeleteByIdRequest) (*DeleteByIdResult, error) {
	err := d.categoryRepo.DeleteById(d.db, request.CategoryId)
	if err != nil {
		d.logger.Error("Failed to delete category by id", slog.String("error", err.Error()))
		return nil, err
	}
	response := &DeleteByIdResult{
		Deleted: true,
		Message: "Category deleted successfully",
	}
	return response, nil
}
