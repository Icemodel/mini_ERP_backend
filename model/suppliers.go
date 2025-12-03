package model

import (
	"time"

	"github.com/google/uuid"
)

type Supplier struct {
	SupplierId uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"supplier_id"`
	Name       string    `gorm:"not null" json:"name"`
	Phone      string    `gorm:"not null" json:"phone"`
	Email      string    `gorm:"not null" json:"email"`
	Address    string    `gorm:"not null" json:"address"`
	CreatedAt  time.Time `gorm:"not null" json:"created_at"`

	PurchaseOrders []PurchaseOrder `gorm:"foreignKey:SupplierId;references:SupplierId"`
}
