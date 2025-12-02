package supplier

import (
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/api/service/supplier/command"
	"mini-erp-backend/api/service/supplier/query"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func NewService(logger *slog.Logger, db *gorm.DB, supplierRepo repository.Supplier) {
	// Register command handlers
	createSupplierHandler := command.NewCreateSupplier(logger, db, supplierRepo)
	updateSupplierHandler := command.NewUpdateSupplier(logger, db, supplierRepo)
	deleteSupplierHandler := command.NewDeleteSupplier(logger, db, supplierRepo)
	
	getSupplierHandler := query.NewGetSupplier(logger, db, supplierRepo)
	getAllSuppliersHandler := query.NewGetAllSuppliers(logger, db, supplierRepo)
	
	err := mediatr.RegisterRequestHandler(createSupplierHandler)
	if err != nil {
		panic(err)
	}

	err = mediatr.RegisterRequestHandler(updateSupplierHandler)
	if err != nil {
		panic(err)
	}

	err = mediatr.RegisterRequestHandler(deleteSupplierHandler)
	if err != nil {
		panic(err)
	}

	// Register query handlers
	err = mediatr.RegisterRequestHandler(getSupplierHandler)
	if err != nil {
		panic(err)
	}

	err = mediatr.RegisterRequestHandler(getAllSuppliersHandler)
	if err != nil {
		panic(err)
	}
}
