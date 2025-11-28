package main

import (
	"fmt"
	"mini-erp-backend/api"
	"mini-erp-backend/api/service/auth"
	"mini-erp-backend/api/service/register"
	"mini-erp-backend/config/database"
	"mini-erp-backend/config/environment"
	"mini-erp-backend/lib/jwt"
	"mini-erp-backend/lib/logging"
	"mini-erp-backend/model"
	"mini-erp-backend/repository"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	log := logging.New()
	environment.LoadEnvironment()

	jwtManager := jwt.New(log.Slogger)

	db := database.Connect(environment.GetString("DSN_DATABASE"))

	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

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

	//region repository
	userAuthen := repository.NewUserAuthen(log.Slogger)
	userRegister := repository.NewUserRegister(log.Slogger)

	//region service
	auth.NewService(db, log.Slogger, jwtManager, userAuthen)
	register.NewService(db, log.Slogger, jwtManager, userRegister)

	app.Listen(":" + environment.GetString("PORT"))
}
