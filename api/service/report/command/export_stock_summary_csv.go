package command

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"log/slog"
	"mini-erp-backend/api/repository"
	"time"

	"gorm.io/gorm"
)

type ExportStockSummaryCSV struct {
	logger     *slog.Logger
	db         *gorm.DB
	reportRepo repository.Report
}

type ExportStockSummaryCSVRequest struct{}

type ExportStockSummaryCSVResult struct {
	Data     []byte
	Filename string
}

func NewExportStockSummaryCSVHandler(
	logger *slog.Logger,
	db *gorm.DB,
	reportRepo repository.Report,
) *ExportStockSummaryCSV {
	return &ExportStockSummaryCSV{
		logger:     logger,
		db:         db,
		reportRepo: reportRepo,
	}
}

func (h *ExportStockSummaryCSV) Handle(ctx context.Context, req *ExportStockSummaryCSVRequest) (*ExportStockSummaryCSVResult, error) {
	products, err := h.reportRepo.GetStockSummary(h.db)
	if err != nil {
		h.logger.Error("Failed to get stock summary for export", "error", err)
		return nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	header := []string{
		"Product Code",
		"Product Name",
		"Category",
		"Stock On Hand",
		"Cost Price",
		"Selling Price",
		"Total Cost Value",
		"Total Selling Value",
		"Min Stock",
		"Low Stock",
	}
	if err := writer.Write(header); err != nil {
		return nil, err
	}

	// Write data rows
	for _, p := range products {
		lowStock := "No"
		if p.StockOnHand < p.MinStock {
			lowStock = "Yes"
		}

		row := []string{
			p.ProductCode,
			p.Name,
			p.CategoryName,
			fmt.Sprintf("%d", p.StockOnHand),
			fmt.Sprintf("%.2f", p.CostPrice),
			fmt.Sprintf("%.2f", p.SellingPrice),
			fmt.Sprintf("%.2f", p.TotalCostValue),
			fmt.Sprintf("%.2f", p.TotalSellingValue),
			fmt.Sprintf("%d", p.MinStock),
			lowStock,
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("stock_summary_%s.csv", time.Now().Format("02-01-2006"))

	return &ExportStockSummaryCSVResult{
		Data:     buf.Bytes(),
		Filename: filename,
	}, nil
}
