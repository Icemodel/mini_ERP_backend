package command

import (
	"context"
	"errors"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StockAdjust struct {
	logger               *slog.Logger
	db                   *gorm.DB
	stockTransactionRepo repository.StockTransaction
	productRepo          repository.Product
}

type StockAdjustRequest struct {
	ProductId uuid.UUID `json:"product_id"`
	Quantity  int64     `json:"quantity"`             // + เพิ่ม, - ลด
	Reason    string    `json:"reason"`               // REQUIRED สำหรับ ADJUST
	CreatedBy string    `json:"created_by,omitempty"` // ผู้ทำรายการ
}

type StockAdjustResult struct {
	Transaction  model.StockTransaction `json:"transaction"`
	CurrentStock int64                  `json:"current_stock"`
	Message      string                 `json:"message"`
}

func NewStockAdjust(logger *slog.Logger, db *gorm.DB, stockTransactionRepo repository.StockTransaction, productRepo repository.Product) *StockAdjust {
	return &StockAdjust{
		logger:               logger,
		db:                   db,
		stockTransactionRepo: stockTransactionRepo,
		productRepo:          productRepo,
	}
}

func (s *StockAdjust) Handle(ctx context.Context, request StockAdjustRequest) (*StockAdjustResult, error) {
	// Validate reason (REQUIRED สำหรับ ADJUST)
	if request.Reason == "" {
		s.logger.Error("Reason is required for stock adjustment")
		return nil, errors.New("reason is REQUIRED for ADJUST transaction type")
	}

	if request.Quantity == 0 {
		s.logger.Error("Quantity cannot be 0")
		return nil, errors.New("quantity cannot be 0 for adjustment")
	}

	// เริ่ม transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	condition := map[string]interface{}{
		"product_id": request.ProductId,
	}

	// ดึงข้อมูล product
	_, err := s.productRepo.Search(tx, condition, "")
	if err != nil {
		tx.Rollback()
		s.logger.Error("Product not found", slog.String("error", err.Error()))
		return nil, errors.New("product not found")
	}

	// สร้าง transaction
	transaction := &model.StockTransaction{
		StockTransactionId: uuid.New(),
		ProductId:          request.ProductId,
		Type:               model.TransactionTypeAdjust,
		Quantity:           request.Quantity,
		Reason:             &request.Reason,
		CreatedAt:          time.Now(),
		CreatedBy:          request.CreatedBy,
	}

	// บันทึก stock adjust
	if err := s.stockTransactionRepo.Create(tx, transaction); err != nil {
		tx.Rollback()
		s.logger.Error("Failed to create stock adjust", slog.String("error", err.Error()))
		return nil, err
	}

	// คำนวณ calculated stock (SUM(IN) - SUM(OUT) + SUM(ADJUST))
	totalIn, totalOut, totalAdjust, err := s.stockTransactionRepo.StockSummary(tx, request.ProductId)
	if err != nil {
		tx.Rollback()
		s.logger.Error("Failed to get stock summary", slog.String("error", err.Error()))
		return nil, err
	}
	currentStock := totalIn - totalOut + totalAdjust

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		s.logger.Error("Failed to commit transaction", slog.String("error", err.Error()))
		return nil, err
	}

	adjustType := "increased"
	if request.Quantity < 0 {
		adjustType = "decreased"
	}

	response := &StockAdjustResult{
		Transaction:  *transaction,
		CurrentStock: currentStock,
		Message:      "Stock ADJUSTED successfully (" + adjustType + ")",
	}

	return response, nil
}
