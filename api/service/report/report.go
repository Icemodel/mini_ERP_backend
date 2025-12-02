package report

import (
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/api/service/report/command"
	"mini-erp-backend/api/service/report/query"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func NewService(logger *slog.Logger, db *gorm.DB, reportRepo repository.Report) error {

	// Register query handlers
	getStockSummaryHandler := query.NewStockSummary(logger, db, reportRepo)
	getStockMovementsHandler := query.NewStockMovements(logger, db, reportRepo)
	getPurchaseSummaryHandler := query.NewPurchaseSummary(logger, db, reportRepo)

	err := mediatr.RegisterRequestHandler(getStockSummaryHandler)
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler(getStockMovementsHandler)
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler(getPurchaseSummaryHandler)
	if err != nil {
		return err
	}

	// Register export handlers
	exportStockSummaryCSVHandler := command.NewExportStockSummaryCSV(logger, db, reportRepo)
	exportStockMovementExcelHandler := command.NewExportStockMovementExcel(logger, db, reportRepo)
	exportPurchaseReportExcelHandler := command.NewExportPurchaseReportExcel(logger, db, reportRepo)

	err = mediatr.RegisterRequestHandler(exportStockSummaryCSVHandler)
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler(exportStockMovementExcelHandler)
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler(exportPurchaseReportExcelHandler)
	if err != nil {
		return err
	}

	logger.Info("Report handlers registered successfully")
	return nil
}
