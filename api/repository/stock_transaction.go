package repository

import (
	"log/slog"
	"mini-erp-backend/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StockTransaction interface {
	// Get
	StockSummary(db *gorm.DB, productId uuid.UUID) (int64, int64, int64, error)
	// Create
	Create(tx *gorm.DB, transaction *model.StockTransaction) error
}

type stockTransaction struct {
	logger *slog.Logger
}

func NewStockTransaction(logger *slog.Logger) StockTransaction {
	return &stockTransaction{
		logger: logger,
	}
}

func (s *stockTransaction) Create(tx *gorm.DB, transaction *model.StockTransaction) error {
	if err := tx.Create(transaction).Error; err != nil {
		s.logger.Error("Failed to create stock in transaction", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (s *stockTransaction) StockSummary(db *gorm.DB, productId uuid.UUID) (int64, int64, int64, error) {
	var totalIn, totalOut, totalAdjust int64

	// Sum IN transactions
	if err := db.Model(&model.StockTransaction{}).
		Where("product_id = ? AND type = ?", productId, model.TransactionTypeIn).
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&totalIn).Error; err != nil {
		s.logger.Error("Failed to sum IN transactions", slog.String("error", err.Error()))
		return 0, 0, 0, err
	}

	// Sum OUT transactions
	if err := db.Model(&model.StockTransaction{}).
		Where("product_id = ? AND type = ?", productId, model.TransactionTypeOut).
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&totalOut).Error; err != nil {
		s.logger.Error("Failed to sum OUT transactions", slog.String("error", err.Error()))
		return 0, 0, 0, err
	}

	// Sum ADJUST transactions
	if err := db.Model(&model.StockTransaction{}).
		Where("product_id = ? AND type = ?", productId, model.TransactionTypeAdjust).
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&totalAdjust).Error; err != nil {
		s.logger.Error("Failed to sum ADJUST transactions", slog.String("error", err.Error()))
		return 0, 0, 0, err
	}

	return totalIn, totalOut, totalAdjust, nil
}

func (s *stockTransaction) GetTransactionsByProduct(db *gorm.DB, productId uuid.UUID) ([]model.StockTransaction, error) {
	var transactions []model.StockTransaction

	if err := db.Where("product_id = ?", productId).Order("created_at DESC").Find(&transactions).Error; err != nil {
		s.logger.Error("Failed to get transactions by product", slog.String("error", err.Error()))
		return nil, err
	}

	return transactions, nil
}
