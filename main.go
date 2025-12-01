package main

import (
	"fmt"
	"mini-erp-backend/api"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/api/service/category"
	"mini-erp-backend/api/service/product"
	"mini-erp-backend/api/service/purchase_order"
	"mini-erp-backend/api/service/report"
	"mini-erp-backend/api/service/stock_transaction"
	"mini-erp-backend/api/service/supplier"
	"mini-erp-backend/config/database"
	"mini-erp-backend/config/environment"
	"mini-erp-backend/lib/logging"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	log := logging.New()

	app.Use(cors.New())

	environment.LoadEnvironment()

	db := database.Connect(environment.GetString("DSN_DATABASE"))

	fmt.Println(db)

	// region Repository
	categoryRepo := repository.NewCategory(log.Slogger)
	productRepo := repository.NewProduct(log.Slogger)
	stockTransactionRepo := repository.NewStockTransaction(log.Slogger)
	supplierRepo := repository.NewSupplier(log.Slogger)
	// endregion

	// region Service
	category.NewService(log.Slogger, db, categoryRepo)
	product.NewService(log.Slogger, db, productRepo, stockTransactionRepo)
	stock_transaction.NewService(log.Slogger, db, stockTransactionRepo, productRepo)
	supplier.RegisterSupplierHandler(log.Slogger, db, supplierRepo)
	if err := purchase_order.RegisterPurchaseOrderHandler(db, log.Slogger); err != nil {
		log.Slogger.Error("Failed to register purchase order handlers", "error", err)
	}
	if err := report.RegisterReportHandlers(log.Slogger, db); err != nil {
		log.Slogger.Error("Failed to register report handlers", "error", err)
	}
	// endregion

	// region Migrations
	// if err := db.AutoMigrate(
	// 	&model.User{},
	// 	&model.Category{},
	// 	&model.Supplier{},
	// 	&model.Product{},
	// 	&model.PurchaseOrder{},
	// 	&model.AuditLog{},
	// 	&model.PurchaseOrderItem{},
	// 	&model.StockTransaction{},
	// ); err != nil {
	// 	log.Slogger.Error("Migration failed", "error", err)
	// }
	// endregion

	// region Routes
	api.Register(app, log.Slogger)
	// endregion

	app.Listen(":" + environment.GetString("PORT"))
}
