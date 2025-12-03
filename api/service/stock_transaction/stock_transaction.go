package stock_transaction

import (
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/api/service/stock_transaction/command"
	"mini-erp-backend/api/service/stock_transaction/query"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func NewService(logger *slog.Logger, db *gorm.DB, stockTransactionRepo repository.StockTransaction, productRepo repository.Product) {
	stockService := query.NewStocks(logger, db, stockTransactionRepo)
	stockInService := command.NewStockIn(logger, db, stockTransactionRepo, productRepo)
	stockOutService := command.NewStockOut(logger, db, stockTransactionRepo, productRepo)
	stockAdjustService := command.NewStockAdjust(logger, db, stockTransactionRepo, productRepo)

	err := mediatr.RegisterRequestHandler(stockService)
	if err != nil {
		panic(err)
	}

	err = mediatr.RegisterRequestHandler(stockInService)
	if err != nil {
		panic(err)
	}

	err = mediatr.RegisterRequestHandler(stockOutService)
	if err != nil {
		panic(err)
	}

	err = mediatr.RegisterRequestHandler(stockAdjustService)
	if err != nil {
		panic(err)
	}
}
