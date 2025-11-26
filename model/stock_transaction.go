package model

import (
	"time"

	"github.com/google/uuid"
)

type StockTransaction struct {
	StockTransactionId uuid.UUID  `gorm:"type:uuid;primaryKey" json:"stock_transaction_id"`
	ProductId          uuid.UUID  `gorm:"type:uuid;not null" json:"product_id"`
	Quantity           int64      `json:"quantity"`
	Type               string     `json:"transaction_type"` // e.g., "IN" or "OUT" or "ADJUST"
	Reason             *string    `json:"reason"`
	ReferenceId        *uuid.UUID `gorm:"type:uuid" json:"reference_id"`
	CreatedAt          time.Time  `json:"created_at"`
	CreatedBy          uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`

	Product Product `gorm:"constraint:OnDelete:CASCADE;" json:"-"`

	// PurchaseOrder *PurchaseOrder `gorm:"foreignKey:ReferenceId"`
	// User User `gorm:"foreignKey:CreatedBy"`
}
