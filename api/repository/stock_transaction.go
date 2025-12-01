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
	FindByProductAndReference(db *gorm.DB, productId uuid.UUID, referenceId uuid.UUID) (*model.StockTransaction, error)
	GetLatestByProductAndReference(db *gorm.DB, productId uuid.UUID, referenceId uuid.UUID) (*model.StockTransaction, error)
	UpdateQuantity(tx *gorm.DB, stockTransactionId uuid.UUID, quantity int64) error
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

func (r *stockTransaction) FindByProductAndReference(db *gorm.DB, productId uuid.UUID, referenceId uuid.UUID) (*model.StockTransaction, error) {
	var stockTx model.StockTransaction
	err := db.Where("product_id = ? AND reference_id = ?", productId, referenceId).First(&stockTx).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // ไม่มีก็ return nil
		}
		r.logger.Error("Failed to find stock transaction", "error", err)
		return nil, err
	}
	return &stockTx, nil
}

func (r *stockTransaction) GetLatestByProductAndReference(db *gorm.DB, productId uuid.UUID, referenceId uuid.UUID) (*model.StockTransaction, error) {
	var stockTx model.StockTransaction
	err := db.Where("product_id = ? AND reference_id = ?", productId, referenceId).
		Order("created_at DESC").
		First(&stockTx).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error("Failed to find latest stock transaction", "error", err)
		return nil, err
	}
	return &stockTx, nil
}

func (r *stockTransaction) UpdateQuantity(tx *gorm.DB, stockTransactionId uuid.UUID, quantity int64) error {
	if err := tx.Model(&model.StockTransaction{}).
		Where("stock_transaction_id = ?", stockTransactionId).
		UpdateColumn("quantity", gorm.Expr("quantity + ?", quantity)).Error; err != nil {
		r.logger.Error("Failed to update stock transaction quantity", "error", err)
		return err
	}
	return nil
}
