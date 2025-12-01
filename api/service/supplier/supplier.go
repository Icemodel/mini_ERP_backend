package supplier

import (
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/api/service/supplier/command"
	"mini-erp-backend/api/service/supplier/query"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func RegisterSupplierHandler(logger *slog.Logger, db *gorm.DB, supplierRepo repository.SupplierRepository) {
	// Register command handlers
	createSupplierHandler := command.NewCreateSupplierHandler(logger, db, supplierRepo)
	updateSupplierHandler := command.NewUpdateSupplierHandler(logger, db, supplierRepo)
	deleteSupplierHandler := command.NewDeleteSupplierHandler(logger, db, supplierRepo)
	getSupplierHandler := query.NewGetSupplierHandler(logger, db, supplierRepo)
	getAllSuppliersHandler := query.NewGetAllSuppliersHandler(logger, db, supplierRepo)
	searchSuppliersHandler := query.NewSearchSuppliersHandler(logger, db, supplierRepo)
	
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

	err = mediatr.RegisterRequestHandler(searchSuppliersHandler)
	if err != nil {
		panic(err)
	}
}
