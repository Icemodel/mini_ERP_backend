package model

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionTypeIn     TransactionType = "IN"
	TransactionTypeOut    TransactionType = "OUT"
	TransactionTypeAdjust TransactionType = "ADJUST"
)

type StockTransaction struct {
	StockTransactionId uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"stock_transaction_id"`
	ProductId          uuid.UUID       `gorm:"type:uuid;not null" json:"product_id"`
	Quantity           int64           `gorm:"not null" json:"quantity"`
	Type               TransactionType `gorm:"not null" json:"type"` // e.g., "IN" or "OUT" or "ADJUST"
	Reason             *string         `json:"reason"`
	ReferenceId        *uuid.UUID      `gorm:"type:uuid" json:"reference_id"`
	CreatedAt          time.Time       `gorm:"not null" json:"created_at"`
	CreatedBy          string          `gorm:"not null" json:"created_by"`

	Product Product `gorm:"constraint:OnDelete:CASCADE;" json:"Product"`

	// PurchaseOrder *PurchaseOrder `gorm:"foreignKey:ReferenceId"`
	// User User `gorm:"foreignKey:CreatedBy"`
}
