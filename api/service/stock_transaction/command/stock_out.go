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

type StockOut struct {
	logger               *slog.Logger
	db                   *gorm.DB
	stockTransactionRepo repository.StockTransaction
	productRepo          repository.Product
}

type StockOutRequest struct {
	ProductId uuid.UUID `json:"product_id"`
	Quantity  int64     `json:"quantity"`
	Reason    *string   `json:"reason"`
	CreatedBy string    `json:"created_by,omitempty"` // ผู้ทำรายการ
}

type StockOutResult struct {
	Transaction  model.StockTransaction `json:"transaction"`
	CurrentStock int64                  `json:"current_stock"`
	Message      string                 `json:"message"`
}

func NewStockOut(logger *slog.Logger, db *gorm.DB, stockTransactionRepo repository.StockTransaction, productRepo repository.Product) *StockOut {
	return &StockOut{
		logger:               logger,
		db:                   db,
		stockTransactionRepo: stockTransactionRepo,
		productRepo:          productRepo,
	}
}

func (s *StockOut) Handle(ctx context.Context, request StockOutRequest) (*StockOutResult, error) {
	// Validate quantity
	if request.Quantity <= 0 {
		s.logger.Error("Quantity must be greater than 0")
		return nil, errors.New("quantity must be greater than 0")
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
		Type:               model.TransactionTypeOut,
		Quantity:           request.Quantity,
		Reason:             request.Reason,
		CreatedAt:          time.Now(),
		CreatedBy:          request.CreatedBy,
	}

	// บันทึก stock out
	if err := s.stockTransactionRepo.Create(tx, transaction); err != nil {
		tx.Rollback()
		s.logger.Error("Failed to create stock out", slog.String("error", err.Error()))
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

	response := &StockOutResult{
		Transaction:  *transaction,
		CurrentStock: currentStock,
		Message:      "Stock OUT successful",
	}

	return response, nil
}