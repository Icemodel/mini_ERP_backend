package main

import (
	"fmt"
	"mini-erp-backend/api"
	"mini-erp-backend/api/service/auth"
	"mini-erp-backend/api/service/category"
	"mini-erp-backend/api/service/product"
	"mini-erp-backend/api/service/purchase_order"
	"mini-erp-backend/api/service/purchase_order_item"
	"mini-erp-backend/api/service/register"
	"mini-erp-backend/api/service/purchase_order_item"
	"mini-erp-backend/api/service/report"
	"mini-erp-backend/api/service/stock_transaction"
	"mini-erp-backend/api/service/supplier"
	"mini-erp-backend/config/database"
	"mini-erp-backend/config/environment"
	"mini-erp-backend/lib/jwt"
	"mini-erp-backend/lib/logging"

	"mini-erp-backend/api/repository"
	_ "mini-erp-backend/docs"
	"mini-erp-backend/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

// Main function
//
//	@title						Mini ERP Backend API
//	@version					1.0
//	@description				This is Mini ERP Backend API doc
//	@termsOfService				http://swagger.io/terms/
//	@BasePath					/api/v1
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Example of "Value": Bearer <your_token>
//	@securityDefinitions.basic	BasicAuth
func main() {
	app := fiber.New()
	log := logging.New()
	environment.LoadEnvironment()

	app.Use(cors.New())
	jwtManager := jwt.New(log.Slogger)

	db := database.Connect(environment.GetString("DSN_DATABASE"))

	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	fmt.Println(db)

	// region Repository
	categoryRepo := repository.NewCategory(log.Slogger)
	productRepo := repository.NewProduct(log.Slogger)
	stockTransactionRepo := repository.NewStockTransaction(log.Slogger)
	supplierRepo := repository.NewSupplier(log.Slogger)
	purchase_orderRepo := repository.NewPurchaseOrder(log.Slogger)
	purchase_order_itemRepo := repository.NewPurchaseOrderItem(log.Slogger)
	reportRepo := repository.NewReport(log.Slogger)
	userRepo := repository.NewUser(log.Slogger)
	sessionRepo := repository.NewUserSession(log.Slogger)
	// endregion

	// region Service
	category.NewService(log.Slogger, db, categoryRepo)
	product.NewService(log.Slogger, db, productRepo, stockTransactionRepo)
	stock_transaction.NewService(log.Slogger, db, stockTransactionRepo, productRepo)
	purchase_order.NewService(db, log.Slogger, purchase_orderRepo)
	purchase_order_item.NewService(db, log.Slogger, purchase_order_itemRepo)
	supplier.NewService(log.Slogger, db, supplierRepo)
	report.NewService(log.Slogger, db, reportRepo)
	auth.NewService(db, log.Slogger, jwtManager, userRepo)
	register.NewService(db, log.Slogger, jwtManager, userRepo)

	// endregion

	// if err := db.AutoMigrate(
		//&model.User{},
		//&model.Category{},
		//&model.Supplier{},
		//&model.Product{},
		//&model.PurchaseOrder{},
		//&model.AuditLog{},
		//&model.PurchaseOrderItem{},
		//&model.StockTransaction{},
	//&model.UserSession{},
	// ); err != nil {
	// 	log.Slogger.Error("Migration failed", "error", err)
	// }

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

	// endregion

	if environment.GetString("ENV") == "development" {
		app.Get("/swagger/*", swagger.HandlerDefault)

		app.Listen(":" + environment.GetString("PORT"))
	}

	//region service

	app.Listen(":" + environment.GetString("PORT"))
}