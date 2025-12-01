package purchase_order

import (
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/api/service/purchase_order/command"
	"mini-erp-backend/api/service/purchase_order/query"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func RegisterPurchaseOrderHandler(db *gorm.DB, logger *slog.Logger) error {
	// Initialize repositories
	poRepo := repository.NewPurchaseOrderRepository(logger)
	stockRepo := repository.NewStockTransactionRepository(logger)
	productRepo := repository.NewProductRepository(logger)

	// Register command handlers
	createPurchaseOrderHandler := command.NewCreatePurchaseOrderHandler(logger, db, poRepo, productRepo)
	updatePurchaseOrderHandler := command.NewUpdatePurchaseOrderHandler(logger, db, poRepo, productRepo)
	updatePOStatusHandler := command.NewUpdatePOStatusHandler(logger, db, poRepo, stockRepo)
	getPurchaseOrderHandler := query.NewGetPurchaseOrderHandler(logger, db, poRepo)
	getAllPurchaseOrdersHandler := query.NewGetAllPurchaseOrdersHandler(logger, db, poRepo)

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
	err = mediatr.RegisterRequestHandler[*query.GetPurchaseOrderRequest, *query.GetPurchaseOrderResult](getPurchaseOrderHandler)
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*query.GetAllPurchaseOrdersRequest, *query.GetAllPurchaseOrdersResult](getAllPurchaseOrdersHandler)
	if err != nil {
		return err
	}

	logger.Info("Purchase Order handlers registered successfully")
	return nil
}
