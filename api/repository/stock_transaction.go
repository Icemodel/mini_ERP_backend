package repository

import (
	"log/slog"
	"mini-erp-backend/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StockTransactionSearchFilters struct {
	Search    string
	ProductId *uuid.UUID
}

type StockTransaction interface {
	// Get
	SearchWithFilters(db *gorm.DB, filters StockTransactionSearchFilters, orderBy string) ([]model.StockTransaction, error)
	SearchWithFiltersAndPagination(db *gorm.DB, filters StockTransactionSearchFilters, orderBy string, page int, pageSize int) ([]model.StockTransaction, int64, error)
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

func (s *stockTransaction) SearchWithFilters(db *gorm.DB, filters StockTransactionSearchFilters, orderBy string) ([]model.StockTransaction, error) {
	var transactions []model.StockTransaction

	query := db.Model(&model.StockTransaction{})

	if filters.Search != "" {
		searchPattern := "%" + filters.Search + "%"
		query = query.Where("quantity ILIKE ? OR type ILIKE ? OR reason ILIKE ?", searchPattern, searchPattern, searchPattern)
	}

	if filters.ProductId != nil {
		query = query.Where("product_id = ?", *filters.ProductId)
	}

	if orderBy != "" {
		query = query.Order(orderBy)
	}

	if err := query.Preload("Product").Preload("Product.Category").Find(&transactions).Error; err != nil {
		s.logger.Error("Failed to search stock transactions with filters", slog.String("error", err.Error()))
		return nil, err
	}

	return transactions, nil
}

func (s *stockTransaction) SearchWithFiltersAndPagination(db *gorm.DB, filters StockTransactionSearchFilters, orderBy string, page int, pageSize int) ([]model.StockTransaction, int64, error) {
	var transactions []model.StockTransaction
	var total int64

	query := db.Model(&model.StockTransaction{})

	if filters.Search != "" {
		searchPattern := "%" + filters.Search + "%"
		query = query.Where("quantity ILIKE ? OR type ILIKE ? OR reason ILIKE ?", searchPattern, searchPattern, searchPattern)
	}
	if filters.ProductId != nil {
		query = query.Where("product_id = ?", *filters.ProductId)
	}

	// นับจำนวนทั้งหมดก่อน pagination
	if err := query.Count(&total).Error; err != nil {
		s.logger.Error("Failed to count stock transactions", slog.String("error", err.Error()))
		return nil, 0, err
	}
	// คำนวณ offset
	offset := (page - 1) * pageSize

	// เรียงลำดับ
	if orderBy != "" {
		query = query.Order(orderBy)
	}

	// ดึงข้อมูลแบบ pagination พร้อม preload Product
	if err := query.Preload("Product").Preload("Product.Category").Limit(pageSize).Offset(offset).Find(&transactions).Error; err != nil {
		s.logger.Error("Failed to search stock transactions with filters and pagination", slog.String("error", err.Error()))
		return nil, 0, err
	}

	return transactions, total, nil
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

func (s *stockTransaction) TransactionsByProduct(db *gorm.DB, productId uuid.UUID) ([]model.StockTransaction, error) {
	var transactions []model.StockTransaction

	if err := db.Where("product_id = ?", productId).Order("created_at DESC").Find(&transactions).Error; err != nil {
		s.logger.Error("Failed to get transactions by product", slog.String("error", err.Error()))
		return nil, err
	}

	return transactions, nil
}
