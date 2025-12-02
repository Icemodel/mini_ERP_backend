package api

import (
	"log/slog"
	auth_handler "mini-erp-backend/api/handler/auth"
	category_handler "mini-erp-backend/api/handler/category"
	product_handler "mini-erp-backend/api/handler/product"
	"mini-erp-backend/api/handler/purchase_order"
	register_handler "mini-erp-backend/api/handler/register"
	"mini-erp-backend/api/handler/report"
	stocktransaction_handler "mini-erp-backend/api/handler/stock_transaction"
	"mini-erp-backend/api/handler/supplier"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, logger *slog.Logger) {
	v1 := app.Group("/api/v1")

	// Supplier routes
	supplierGroup := v1.Group("/suppliers")
	{
		supplierGroup.Get("/", supplier.AllSuppliers(logger))
		supplierGroup.Post("/", supplier.CreateSupplier(logger))
		supplierGroup.Get("/:id", supplier.Supplier(logger))
		supplierGroup.Put("/:id", supplier.UpdateSupplier(logger))
		supplierGroup.Delete("/:id", supplier.DeleteSupplier(logger))
	}

	// Purchase Order routes
	purchaseOrderGroup := v1.Group("/purchase-orders")
	{
		purchaseOrderGroup.Get("/", purchase_order.AllPurchaseOrders(logger))
		purchaseOrderGroup.Post("/", purchase_order.CreatePurchaseOrder(logger))
		purchaseOrderGroup.Get("/:id", purchase_order.PurchaseOrder(logger))
		purchaseOrderGroup.Put("/:id", purchase_order.UpdatePurchaseOrder(logger))
		purchaseOrderGroup.Put("/:id/status", purchase_order.UpdatePurchaseOrderStatus(logger))
	}

	// Report routes
	reportGroup := v1.Group("/reports")
	{
		reportGroup.Get("/stock-summary", report.StockSummary(logger))
		reportGroup.Get("/stock-summary/export", report.ExportStockSummaryCSV(logger))
		reportGroup.Get("/stock-movements", report.StockMovements(logger))
		reportGroup.Get("/stock-movements/export", report.ExportStockMovementExcel(logger))
		reportGroup.Get("/purchase-summary", report.PurchaseSummary(logger))
		reportGroup.Get("/purchase-summary/export", report.ExportPurchaseReportExcel(logger))
	}

	categoryGroupApi := v1.Group("/categories")
	{
		categoryGroupApi.Get("/", category_handler.Categories(logger))
		categoryGroupApi.Get("/:id", category_handler.CategoryById(logger))
		categoryGroupApi.Post("/", category_handler.Create(logger))
		categoryGroupApi.Patch("/:id", category_handler.Update(logger))
		categoryGroupApi.Delete("/:id", category_handler.DeleteById(logger))
	}

	productGroupApi := v1.Group("/products")
	{
		productGroupApi.Get("/", product_handler.Products(logger))
		productGroupApi.Get("/:id", product_handler.ProductById(logger))
		productGroupApi.Post("/", product_handler.Create(logger))
		productGroupApi.Patch("/:id", product_handler.Update(logger))
		productGroupApi.Delete("/:id", product_handler.DeleteById(logger))
		productGroupApi.Get("/:id/stock-summary", product_handler.ProductStockSummary(logger))
	}

	stockGroupApi := v1.Group("/stocks")
	{
		stockGroupApi.Get("/", stocktransaction_handler.StockTransactions(logger))
		stockGroupApi.Post("/in", stocktransaction_handler.StockIn(logger))
		stockGroupApi.Post("/out", stocktransaction_handler.StockOut(logger))
		stockGroupApi.Post("/adjust", stocktransaction_handler.StockAdjust(logger))
	}

	authGroupApi := v1.Group("/auth")
	{
		authGroupApi.Post("/login", auth_handler.Login(logger))
		authGroupApi.Post("/register", register_handler.Register(logger))
	}
}
