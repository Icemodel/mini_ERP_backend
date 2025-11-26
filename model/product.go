package model

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ProductId    uuid.UUID `gorm:"type:uuid;primaryKey" json:"product_id"`
	ProductCode  string    `json:"product_code"`
	CategoryId   uuid.UUID `gorm:"type:uuid;not null" json:"category_id"`
	Name         string    `json:"name"`
	CostPrice    float64   `json:"cost_price"`
	SellingPrice float64   `json:"selling_price"`
	Unit         int64     `json:"unit"`
	MinStock     int64     `json:"min_stock"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Category          Category           `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	StockTransactions []StockTransaction `gorm:"foreignKey:ProductId;references:ProductId"`
}
