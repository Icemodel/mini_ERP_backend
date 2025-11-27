package model

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ProductId    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"product_id"`
	ProductCode  string    `gorm:"not null;uniqueIndex" json:"product_code"`
	CategoryId   uuid.UUID `gorm:"type:uuid;not null" json:"category_id"`
	Name         string    `gorm:"not null" json:"name"`
	CostPrice    float64   `gorm:"not null" json:"cost_price"`
	SellingPrice float64   `gorm:"not null" json:"selling_price"`
	Unit         int64     `gorm:"not null" json:"unit"`
	MinStock     int64     `gorm:"not null" json:"min_stock"`
	CreatedAt    time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null" json:"updated_at"`

	Category           Category            `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	StockTransactions  []StockTransaction  `gorm:"foreignKey:ProductId;references:ProductId"`
	PurchaseOrderItems []PurchaseOrderItem `gorm:"foreignKey:ProductId;references:ProductId" json:"-"`
}
