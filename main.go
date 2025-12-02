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
	"mini-erp-backend/middleware"
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
	//&model.User{},
	//&model.Category{},
	//&model.Supplier{},
	//&model.Product{},
	//&model.PurchaseOrder{},
	//&model.AuditLog{},
	//&model.PurchaseOrderItem{},
	//&model.StockTransaction{},
	//&model.UserSession{},
	); err != nil {
		log.Slogger.Error("Migration failed", "error", err)
	}

	//region repository
	userRepo := repository.NewUser(log.Slogger)

	// session repository (user sessions: access/refresh tokens)
	sessionRepo := repository.NewUserSession(log.Slogger)

	//region service
	auth.NewService(db, log.Slogger, jwtManager, userRepo)
	register.NewService(db, log.Slogger, jwtManager, userRepo)

	//middleware
	mid := middleware.NewFiberMiddleware(
		db,
		log.Slogger,
		jwtManager,
		userRepo,
		sessionRepo,
	)
	app.Use(mid.CORS())

	// region Routes
	api.Register(
		app,
		log.Slogger,
		jwtManager,
		mid,
	)

	app.Listen(":" + environment.GetString("PORT"))
}
