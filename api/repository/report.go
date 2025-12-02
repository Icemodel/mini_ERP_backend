package repository

import (
	"log/slog"
	"mini-erp-backend/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Report interface {
	GetStockSummary(db *gorm.DB) ([]StockSummaryResult, error)
	GetStockMovements(db *gorm.DB, fromDate, toDate time.Time) ([]StockMovementResult, error)
	GetPurchaseSummary(db *gorm.DB, year int, month int) ([]PurchaseSummaryResult, error)
}

type report struct{
	logger *slog.Logger
}

func NewReport(logger *slog.Logger) Report {
	return &report{
		logger: logger,
	}
}

// Result structs
type StockSummaryResult struct {
	ProductId         uuid.UUID `json:"product_id"`
	ProductCode       string    `json:"product_code"`
	Name              string    `json:"name"`
	StockOnHand       int64     `json:"stock_on_hand"`
	CostPrice         float64   `json:"cost_price"`
	SellingPrice      float64   `json:"selling_price"`
	TotalCostValue    float64   `json:"total_cost_value"`
	TotalSellingValue float64   `json:"total_selling_value"`
	MinStock          int64     `json:"min_stock"`
	CategoryName      string    `json:"category_name"`
	IsLowStock        bool      `json:"is_low_stock"`
}

type StockMovementResult struct {
	StockTransactionId uuid.UUID  `json:"stock_transaction_id"`
	ProductId          uuid.UUID  `json:"product_id"`
	ProductCode        string     `json:"product_code"`
	ProductName        string     `json:"product_name"`
	CategoryName       string     `json:"category_name"`
	Quantity           int64      `json:"quantity"`
	Type               string     `json:"type"`
	Reason             *string    `json:"reason"`
	ReferenceId        *uuid.UUID `json:"reference_id"`
	CreatedAt          time.Time  `json:"created_at"`
	CreatedBy          string     `json:"created_by"`
}

type PurchaseSummaryResult struct {
	Status        string  `json:"status"`
	TotalOrders   int64   `json:"total_orders"`
	TotalAmount   uint64  `json:"total_amount"`
	AverageAmount float64 `json:"average_amount"`
}

// GetStockSummary returns stock summary with cost and selling values
func (r *report) GetStockSummary(db *gorm.DB) ([]StockSummaryResult, error) {
	var results []StockSummaryResult

	err := db.Table("products").
		Select(`
			products.product_id,
			products.product_code,
			products.name,
			COALESCE(latest_stock.quantity, 0) as stock_on_hand,
			products.cost_price,
			products.selling_price,
			(COALESCE(latest_stock.quantity, 0) * products.cost_price) as total_cost_value,
			(COALESCE(latest_stock.quantity, 0) * products.selling_price) as total_selling_value,
			products.min_stock,
			categories.name as category_name
		`).
		Joins("LEFT JOIN categories ON products.category_id = categories.category_id").
		Joins(`LEFT JOIN LATERAL (
			SELECT quantity 
			FROM stock_transactions 
			WHERE stock_transactions.product_id = products.product_id 
			ORDER BY created_at DESC 
			LIMIT 1
		) latest_stock ON true`).
		Order("products.name ASC").
		Scan(&results).Error

	return results, err
}

// GetStockMovements returns stock transactions within a date range
func (r *report) GetStockMovements(db *gorm.DB, fromDate, toDate time.Time) ([]StockMovementResult, error) {
	var results []StockMovementResult

	err := db.Model(&model.StockTransaction{}).
		Select(`
			stock_transactions.stock_transaction_id,
			stock_transactions.product_id,
			stock_transactions.quantity,
			stock_transactions.type,
			stock_transactions.reason,
			stock_transactions.reference_id,
			stock_transactions.created_at,
			stock_transactions.created_by,
			products.product_code,
			products.name as product_name,
			categories.name as category_name
		`).
		Joins("LEFT JOIN products ON stock_transactions.product_id = products.product_id").
		Joins("LEFT JOIN categories ON products.category_id = categories.category_id").
		Where("stock_transactions.created_at >= ? AND stock_transactions.created_at <= ?", fromDate, toDate).
		Order("stock_transactions.created_at DESC").
		Scan(&results).Error

	return results, err
}

// GetPurchaseSummary returns monthly purchase order summary grouped by status
func (r *report) GetPurchaseSummary(db *gorm.DB, year int, month int) ([]PurchaseSummaryResult, error) {
	var results []PurchaseSummaryResult

	// Calculate first and last day of month
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1).Add(24*time.Hour - time.Second)

	err := db.Model(&model.PurchaseOrder{}).
		Select(`
			status,
			COUNT(*) as total_orders,
			SUM(total_amount) as total_amount,
			AVG(total_amount) as average_amount
		`).
		Where("created_at >= ? AND created_at <= ?", firstDay, lastDay).
		Group("status").
		Order("status").
		Scan(&results).Error

	return results, err
}
