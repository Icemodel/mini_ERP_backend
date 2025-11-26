package main

import (
	"fmt"
	"mini-erp-backend/api"
	"mini-erp-backend/config/database"
	"mini-erp-backend/config/environment"
	"mini-erp-backend/lib/logging"
	"mini-erp-backend/model"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	log := logging.New()

	environment.LoadEnvironment()

	db := database.Connect(environment.GetString("DSN_DATABASE"))

	fmt.Println(db)

	// region Repository
	// productRepo := repository.NewProduct(log.Slogger)
	// endregion

	// region Service
	// product.NewService(log.Slogger, db, productRepo)
	// endregion

	// region Migrations
	if err := db.AutoMigrate(
		&model.User{},
    	&model.Category{},
    	&model.Product{},
    	&model.Supplier{},
    	&model.PurchaseOrder{},
    	&model.PurchaseOrderItem{},
    	&model.StockTransaction{},
    	&model.AuditLog{},
	); err != nil {
		log.Slogger.Error("Migration failed", "error", err)
	}
	// endregion

	// region Routes
	api.Register(app, log.Slogger)
	// endregion

	app.Listen(":" + environment.GetString("PORT"))
}
