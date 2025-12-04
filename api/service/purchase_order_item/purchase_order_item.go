package purchase_order_item

import (
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/api/service/purchase_order_item/command"
	"mini-erp-backend/api/service/purchase_order_item/query"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func NewService(db *gorm.DB, logger *slog.Logger,purchaseOrderItemRepo repository.PurchaseOrderItem ) error {
	// Initialize repositories
	poItemRepo := repository.NewPurchaseOrderItem(logger)
	poRepo := repository.NewPurchaseOrder(logger)
	productRepo := repository.NewProduct(logger)

	// Register command handlers
	createHandler := command.NewCreatePurchaseOrderItem(logger, db, poItemRepo, poRepo, productRepo)
	updateHandler := command.NewUpdatePurchaseOrderItem(logger, db, poItemRepo, poRepo)
	deleteHandler := command.NewDeletePurchaseOrderItem(logger, db, poItemRepo, poRepo)

	err := mediatr.RegisterRequestHandler(createHandler)
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler(updateHandler)
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler(deleteHandler)
	if err != nil {
		return err
	}

	// Register query handlers
	itemHandler := query.NewPurchaseOrderItem(logger, db, poItemRepo)
	itemsHandler := query.NewPurchaseOrderItems(logger, db, poItemRepo)
	allItemsHandler := query.NewAllPurchaseOrderItems(logger, db, poItemRepo)

	err = mediatr.RegisterRequestHandler[*query.PurchaseOrderItemRequest, interface{}](itemHandler)
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*query.PurchaseOrderItemsRequest, interface{}](itemsHandler)
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*query.AllPurchaseOrderItemsRequest, interface{}](allItemsHandler)
	if err != nil {
		return err
	}

	logger.Info("Purchase Order Item handlers registered successfully")
	return nil
}
