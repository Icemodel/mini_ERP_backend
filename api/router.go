package api

import (
	"log/slog"
	auth_handler "mini-erp-backend/api/handler/auth"
	category_handler "mini-erp-backend/api/handler/category"
	product_handler "mini-erp-backend/api/handler/product"
	"mini-erp-backend/api/handler/purchase_order"
	register_handler "mini-erp-backend/api/handler/register"
	"mini-erp-backend/api/handler/purchase_order_item"
	"mini-erp-backend/api/handler/report"
	stocktransaction_handler "mini-erp-backend/api/handler/stock_transaction"
	"mini-erp-backend/api/handler/supplier"
	"mini-erp-backend/lib/jwt"
	"mini-erp-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func Register(
	app *fiber.App,
	logger *slog.Logger,
	jwt jwt.Manager,
	mid *middleware.FiberMiddleware,
) {
	v1 := app.Group("/api/v1")

	// Auth routes
	authGroupApi := v1.Group("/auth")
	{
		authGroupApi.Post("/login", auth_handler.Login(logger))
		authGroupApi.Post("/token/refresh", auth_handler.RefreshAccessToken(logger))
		authGroupApi.Post("/logout", auth_handler.Logout(logger))
	}

	// Supplier routes
	supplierGroup := v1.Group("/suppliers")
	{
		supplierGroup.Get("/", mid.RequireMinRole("viewer"), supplier.AllSuppliers(logger))
		supplierGroup.Get("/:id", mid.RequireMinRole("admin"), supplier.Supplier(logger))
		supplierGroup.Post("/", mid.RequireMinRole("admin"), supplier.CreateSupplier(logger))
		supplierGroup.Put("/:id", mid.RequireMinRole("admin"), supplier.UpdateSupplier(logger))
		supplierGroup.Delete("/:id", mid.RequireMinRole("admin"), supplier.DeleteSupplier(logger))
		
	}

	// Purchase Order routes
	purchaseOrderGroup := v1.Group("/purchase-orders")
	{
		purchaseOrderGroup.Use(mid.Authenticated())

		purchaseOrderGroup.Get("/", mid.RequireMinRole("viewer"), purchase_order.AllPurchaseOrders(logger))
		purchaseOrderGroup.Get("/:id", mid.RequireMinRole("staff"), purchase_order.PurchaseOrder(logger))
		purchaseOrderGroup.Post("/", mid.RequireMinRole("staff"), purchase_order.CreatePurchaseOrder(logger))
		purchaseOrderGroup.Put("/:id", mid.RequireMinRole("staff"), purchase_order.UpdatePurchaseOrder(logger))
		purchaseOrderGroup.Put("/:id/status", mid.RequireMinRole("staff"), purchase_order.UpdatePurchaseOrderStatus(logger))

	}

	// Purchase Order Item routes
	purchaseOrderItemGroup := v1.Group("/purchase-order-items")
	{
		purchaseOrderItemGroup.Get("/", purchase_order_item.AllPurchaseOrderItems(logger))
		purchaseOrderItemGroup.Get("/:po_id", purchase_order_item.PurchaseOrderItems(logger))
		purchaseOrderItemGroup.Get("/item/:item_id", purchase_order_item.PurchaseOrderItem(logger))
		purchaseOrderItemGroup.Post("/", purchase_order_item.CreatePurchaseOrderItem(logger))
		purchaseOrderItemGroup.Put("/:item_id", purchase_order_item.UpdatePurchaseOrderItem(logger))
		purchaseOrderItemGroup.Delete("/:item_id", purchase_order_item.DeletePurchaseOrderItem(logger))
	}

	// Report routes
	reportGroup := v1.Group("/reports")
	{
		reportGroup.Use(mid.Authenticated())

		reportGroup.Get("/stock-summary", mid.RequireMinRole("admin"), report.StockSummary(logger))
		reportGroup.Get("/stock-summary/export", mid.RequireMinRole("admin"), report.ExportStockSummaryCSV(logger))
		reportGroup.Get("/stock-movements", mid.RequireMinRole("admin"), report.StockMovements(logger))
		reportGroup.Get("/stock-movements/export", mid.RequireMinRole("admin"), report.ExportStockMovementExcel(logger))
		reportGroup.Get("/purchase-summary", mid.RequireMinRole("admin"), report.PurchaseSummary(logger))
		reportGroup.Get("/purchase-summary/export", mid.RequireMinRole("admin"), report.ExportPurchaseReportExcel(logger))

	}

	categoryGroupApi := v1.Group("/categories")
	{
		categoryGroupApi.Use(mid.Authenticated())

		categoryGroupApi.Get("/", mid.RequireMinRole("viewer"), category_handler.Categories(logger))
		categoryGroupApi.Get("/:id", mid.RequireMinRole("viewer"), category_handler.CategoryById(logger))
		categoryGroupApi.Post("/", mid.RequireMinRole("admin"), category_handler.Create(logger))
		categoryGroupApi.Patch("/:id", mid.RequireMinRole("admin"), category_handler.Update(logger))
		categoryGroupApi.Delete("/:id", mid.RequireMinRole("admin"), category_handler.DeleteById(logger))
	}

	productGroupApi := v1.Group("/products")
	{
		productGroupApi.Use(mid.Authenticated())

		productGroupApi.Get("/", mid.RequireMinRole("viewer"), product_handler.Products(logger))
		productGroupApi.Get("/:id", mid.RequireMinRole("viewer"), product_handler.ProductById(logger))
		productGroupApi.Post("/", mid.RequireMinRole("staff"), product_handler.Create(logger))
		productGroupApi.Patch("/:id", mid.RequireMinRole("staff"), product_handler.Update(logger))
		productGroupApi.Delete("/:id", mid.RequireMinRole("staff"), product_handler.DeleteById(logger))
		productGroupApi.Get("/:id/stock-summary", mid.RequireMinRole("viewer"), product_handler.ProductStockSummary(logger))
	}

	stockGroupApi := v1.Group("/stocks")
	{
		stockGroupApi.Use(mid.Authenticated())

		stockGroupApi.Get("/", mid.RequireMinRole("viewer"), stocktransaction_handler.StockTransactions(logger))
		stockGroupApi.Post("/in", mid.RequireMinRole("staff"), stocktransaction_handler.StockIn(logger))
		stockGroupApi.Post("/out", mid.RequireMinRole("staff"), stocktransaction_handler.StockOut(logger))
		stockGroupApi.Post("/adjust", mid.RequireMinRole("staff"), stocktransaction_handler.StockAdjust(logger))
	}

	//Test Route (Add user regis)
	registerGroupApi := v1.Group("/register") // Test only
	{
		registerGroupApi.Use(mid.Authenticated())

		registerGroupApi.Post("/", mid.RequireRole("admin"), register_handler.Register(logger))
	}
}
