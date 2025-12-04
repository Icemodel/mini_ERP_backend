package model

import (
	"time"

	"github.com/google/uuid"
)

type PurchaseOrder struct {
	PurchaseOrderId uuid.UUID           `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"purchase_order_id"`
	SupplierId      uuid.UUID           `gorm:"type:uuid;not null;" json:"supplier_id"`
	Status          PurchaseOrderStatus `gorm:"not null" json:"status"`
	// TotalAmount     uint64              `gorm:"not null" json:"total_amount"`
	CreatedAt       time.Time           `gorm:"not null" json:"created_at"`
	CreatedBy       string          `gorm:"not null" json:"created_by"`

	PurchaseOrderItem []PurchaseOrderItem `gorm:"foreignKey:PurchaseOrderId" json:"purchase_order_items"`
	StockTransaction  []StockTransaction  `gorm:"foreignKey:ReferenceId;constraint:OnDelete:SET NULL;" json:"stock_transactions"`
	Supplier          Supplier            `gorm:"constraint:OnDelete:SET NULL;" json:"-"`
}

type PurchaseOrderStatus string

const (
	Draft     PurchaseOrderStatus = "DRAFT"
	Confirmed PurchaseOrderStatus = "CONFIRMED"
	Received  PurchaseOrderStatus = "RECEIVED"
	Cancelled PurchaseOrderStatus = "CANCELLED"
)
