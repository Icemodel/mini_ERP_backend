package command

import (
	"context"
	"fmt"
	"log/slog"
	"mini-erp-backend/api/repository"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type ExportPurchaseReportExcel struct {
	logger     *slog.Logger
	db         *gorm.DB
	reportRepo repository.Report
}

type ExportPurchaseReportExcelRequest struct {
	Year  int
	Month int
}

type ExportPurchaseReportExcelResult struct {
	Data     []byte
	Filename string
}

func NewExportPurchaseReportExcel(
	logger *slog.Logger,
	db *gorm.DB,
	reportRepo repository.Report,
) *ExportPurchaseReportExcel {
	return &ExportPurchaseReportExcel{
		logger:     logger,
		db:         db,
		reportRepo: reportRepo,
	}
}

func (h *ExportPurchaseReportExcel) Handle(ctx context.Context, req *ExportPurchaseReportExcelRequest) (*ExportPurchaseReportExcelResult, error) {
	summary, err := h.reportRepo.GetPurchaseSummary(h.db, req.Year, req.Month)
	if err != nil {
		h.logger.Error("Failed to get purchase summary for export", "error", err)
		return nil, err
	}

	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Purchase Report"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Create header style
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 11, Color: "#FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#4472C4"}, Pattern: 1},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, err
	}

	// Write title
	f.SetCellValue(sheetName, "A1", fmt.Sprintf("Purchase Order Summary - %02d/%d", req.Month, req.Year))
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	f.SetCellStyle(sheetName, "A1", "D1", titleStyle)
	f.MergeCell(sheetName, "A1", "D1")

	// Write headers
	headers := []string{"Status", "Total Orders", "Total Amount", "Average Amount"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 3)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Write data
	var totalOrders int64
	var totalAmount uint64
	for i, s := range summary {
		row := i + 4
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), s.Status)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), s.TotalOrders)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), s.TotalAmount)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), fmt.Sprintf("%.2f", s.AverageAmount))

		totalOrders += s.TotalOrders
		totalAmount += s.TotalAmount
	}

	// Write summary row
	summaryRow := len(summary) + 5
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", summaryRow), "TOTAL")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", summaryRow), totalOrders)
	f.SetCellValue(sheetName, fmt.Sprintf("C%d", summaryRow), totalAmount)

	summaryStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 11},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E7E6E6"}, Pattern: 1},
	})
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", summaryRow), fmt.Sprintf("D%d", summaryRow), summaryStyle)

	// Auto-fit columns
	f.SetColWidth(sheetName, "A", "A", 15)
	f.SetColWidth(sheetName, "B", "B", 15)
	f.SetColWidth(sheetName, "C", "C", 18)
	f.SetColWidth(sheetName, "D", "D", 18)

	// Freeze header row
	f.SetPanes(sheetName, &excelize.Panes{
		Freeze:      true,
		XSplit:      0,
		YSplit:      3,
		TopLeftCell: "A4",
		ActivePane:  "bottomLeft",
	})

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("purchase_report_%02d-%d.xlsx", req.Month, req.Year)

	return &ExportPurchaseReportExcelResult{
		Data:     buffer.Bytes(),
		Filename: filename,
	}, nil
}
