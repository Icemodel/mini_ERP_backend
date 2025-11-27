package main

import (
	"fmt"
	"mini-erp-backend/api"
	"mini-erp-backend/config/database"
	"mini-erp-backend/config/environment"
	"mini-erp-backend/lib/jwt"
	"mini-erp-backend/lib/logging"
	"mini-erp-backend/model"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	log := logging.New()
	jwtManager := jwt.New(log.Slogger)

	environment.LoadEnvironment()

	db := database.Connect(environment.GetString("DSN_DATABASE"))

	fmt.Println(db)

	if err := db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Supplier{},
		&model.Product{},
		&model.PurchaseOrder{},
		&model.AuditLog{},
		&model.PurchaseOrderItem{},
		&model.StockTransaction{},
	); err != nil {
		log.Slogger.Error("Migration failed", "error", err)
	}
	// endregion

	// region Routes
	api.Register(
		app,
		log.Slogger,
		jwtManager,
	)

	// endregion

	app.Listen(":" + environment.GetString("PORT"))
}
