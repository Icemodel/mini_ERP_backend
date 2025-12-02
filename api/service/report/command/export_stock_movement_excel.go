package command

import (
	"context"
	"fmt"
	"log/slog"
	"mini-erp-backend/api/repository"
	"time"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type ExportStockMovementExcel struct {
	logger     *slog.Logger
	db         *gorm.DB
	reportRepo repository.Report
}

type ExportStockMovementExcelRequest struct {
	FromDate time.Time
	ToDate   time.Time
}

type ExportStockMovementExcelResult struct {
	Data     []byte
	Filename string
}

func NewExportStockMovementExcel(
	logger *slog.Logger,
	db *gorm.DB,
	reportRepo repository.Report,
) *ExportStockMovementExcel {
	return &ExportStockMovementExcel{
		logger:     logger,
		db:         db,
		reportRepo: reportRepo,
	}
}

func (h *ExportStockMovementExcel) Handle(ctx context.Context, req *ExportStockMovementExcelRequest) (*ExportStockMovementExcelResult, error) {
	movements, err := h.reportRepo.GetStockMovements(h.db, req.FromDate, req.ToDate)
	if err != nil {
		h.logger.Error("Failed to get stock movements for export", "error", err)
		return nil, err
	}

	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Stock Movements"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Create header style
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 11},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#D3D3D3"}, Pattern: 1},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, err
	}

	// Write headers
	headers := []string{"Transaction ID", "Date", "Product Code", "Product Name", "Category", "Type", "Quantity", "Reason", "Reference ID", "Created By"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Write data
	for i, m := range movements {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), m.StockTransactionId.String())
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), m.CreatedAt.Format("02-01-2006 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), m.ProductCode)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), m.ProductName)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), m.CategoryName)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), m.Type)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), m.Quantity)
		
		reason := ""
		if m.Reason != nil {
			reason = *m.Reason
		}
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), reason)

		refId := ""
		if m.ReferenceId != nil {
			refId = m.ReferenceId.String()
		}
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), refId)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), m.CreatedBy)
	}

	// Auto-fit columns
	for i := 1; i <= len(headers); i++ {
		col, _ := excelize.ColumnNumberToName(i)
		f.SetColWidth(sheetName, col, col, 15)
	}

	// Freeze first row
	f.SetPanes(sheetName, &excelize.Panes{
		Freeze:      true,
		XSplit:      0,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("stock_movements_%s_to_%s.xlsx",
		req.FromDate.Format("02-01-2006"),
		req.ToDate.Format("02-01-2006"))

	return &ExportStockMovementExcelResult{
		Data:     buffer.Bytes(),
		Filename: filename,
	}, nil
}
