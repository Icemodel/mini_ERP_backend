package purchase_order

import (
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/api/service/purchase_order/command"
	"mini-erp-backend/api/service/purchase_order/query"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func NewService(db *gorm.DB, logger *slog.Logger, purchaseOrderRepo repository.PurchaseOrder) error {
	// Initialize repositories
	poRepo := repository.NewPurchaseOrder(logger)
	stockRepo := repository.NewStockTransaction(logger)
	productRepo := repository.NewProduct(logger)
	poItemRepo := repository.NewPurchaseOrderItem(logger)

	// Register command handlers
	createPurchaseOrderHandler := command.NewCreatePurchaseOrder(logger, db, poRepo, productRepo)
	updatePurchaseOrderHandler := command.NewUpdatePurchaseOrder(logger, db, poRepo, productRepo)
	updatePOStatusHandler := command.NewUpdatePOStatus(logger, db, poRepo, poItemRepo, stockRepo)
	getPurchaseOrderHandler := query.NewPurchaseOrder(logger, db, poRepo)
	getAllPurchaseOrdersHandler := query.NewAllPurchaseOrders(logger, db, poRepo)

	err := mediatr.RegisterRequestHandler(createPurchaseOrderHandler)
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler(updatePurchaseOrderHandler)
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler(updatePOStatusHandler)
	if err != nil {
		return err
	}

	// Register query handlers
	err = mediatr.RegisterRequestHandler(getPurchaseOrderHandler)
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler(getAllPurchaseOrdersHandler)
	if err != nil {
		return err
	}

	return nil
}
