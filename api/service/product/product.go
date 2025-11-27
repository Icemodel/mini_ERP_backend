package product

import (
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/api/service/product/command"
	"mini-erp-backend/api/service/product/query"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func NewService(logger *slog.Logger, db *gorm.DB, productRepo repository.Product) {
	productService := query.NewProducts(logger, db, productRepo)
	productByIdService := query.NewProductById(logger, db, productRepo)
	createProductService := command.NewCreate(logger, db, productRepo)
	updateProductService := command.NewUpdate(logger, db, productRepo)
	deleteProductByIdService := command.NewDeleteById(logger, db, productRepo)

	err := mediatr.RegisterRequestHandler(productService)
	if err != nil {
		panic(err)
	}

	err = mediatr.RegisterRequestHandler(productByIdService)
	if err != nil {
		panic(err)
	}

	err = mediatr.RegisterRequestHandler(createProductService)
	if err != nil {
		panic(err)
	}

	err = mediatr.RegisterRequestHandler(updateProductService)
	if err != nil {
		panic(err)
	}

	err = mediatr.RegisterRequestHandler(deleteProductByIdService)
	if err != nil {
		panic(err)
	}
}
