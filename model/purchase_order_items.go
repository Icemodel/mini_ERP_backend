package model

import (
	"github.com/google/uuid"
)

type PurchaseOrderItem struct {
	PurchaseOrderItemId uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"purchase_order_item_id"`
	PurchaseOrderId     uuid.UUID `gorm:"type:uuid;not null;" json:"purchase_order_id"`
	ProductId           uuid.UUID `gorm:"type:uuid;not null;" json:"product_id"`
	Quantity            uint64    `gorm:"not null;" json:"quantity"`
	Price               float64   `gorm:"not null;" json:"price"`

	PurchaseOrder PurchaseOrder `gorm:"foconstraint:OnDelete:CASCADE;" json:"-"`
	Product       Product       `gorm:"constraint:OnDelete:SET NULL;" json:"-"`
}
