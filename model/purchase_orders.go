package model

import (
	"time"

	"github.com/google/uuid"
)

type PurchaseOrder struct {
	PurchaseOrderId uuid.UUID           `gorm:"type:uuid;primaryKey" json:"purchase_order_id"`
	SupplierId      uuid.UUID           `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"supplier_id"`
	Supplier        Supplier            `gorm:"foreignKey:SupplierId" json:"-"`
	Items           []PurchaseOrderItem `gorm:"foreignKey:PurchaseOrderId" json:"items"`
	Status          PurchaseOrderStatus `json:"status"`
	TotalAmount     uint64              `json:"total_amount"`
	CreatedAt       time.Time           `json:"created_at"`
	CreatedBy       string              `gorm:"type:uuid" json:"created_by"`
}

type PurchaseOrderStatus string

const (
	Draft     PurchaseOrderStatus = "DRAFT"
	Confirmed PurchaseOrderStatus = "CONFIRMED"
	Received  PurchaseOrderStatus = "RECEIVED"
	Cancelled PurchaseOrderStatus = "CANCELLED"
)
