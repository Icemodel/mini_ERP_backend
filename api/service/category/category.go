package category

import (
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/api/service/category/command"
	"mini-erp-backend/api/service/category/query"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func NewService(logger *slog.Logger, db *gorm.DB, categoryRepo repository.Category) {
	categoryService := query.NewCategories(logger, db, categoryRepo)
	categoryByIdService := query.NewCategoryById(logger, db, categoryRepo)
	createCategoryService := command.NewCreate(logger, db, categoryRepo)
	updateCategoryService := command.NewUpdate(logger, db, categoryRepo)
	deleteCategoryByIdService := command.NewDeleteById(logger, db, categoryRepo)

	err := mediatr.RegisterRequestHandler(categoryService)
	if err != nil {
		panic(err)
	}

	err = mediatr.RegisterRequestHandler(categoryByIdService)
	if err != nil {
		panic(err)
	}

	err = mediatr.RegisterRequestHandler(createCategoryService)
	if err != nil {
		panic(err)
	}

	err = mediatr.RegisterRequestHandler(updateCategoryService)
	if err != nil {
		panic(err)
	}

	err = mediatr.RegisterRequestHandler(deleteCategoryByIdService)
	if err != nil {
		panic(err)
	}
}
