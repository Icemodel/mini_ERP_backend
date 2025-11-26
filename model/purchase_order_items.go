package model

import (
	"github.com/google/uuid"
)

type Purchase_order_item  struct {
	PurchaseOrderItemId   uuid.UUID `gorm:"type:uuid;primaryKey" json:"purchase_order_item_id"`
	PurchaseOrderId  uuid.UUID `gorm:"type:uuid;index;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"purchase_order_id"`
	PurchaseOrder   Purchase_order `gorm:"foreignKey:PurchaseOrderId;references:PurchaseOrderId" json:"purchase_order"`
	ProductId uuid.UUID `gorm:"type:uuid;index;not null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductId;references:ProductId" json:"product"`
	Quantity uint64 `json:"quantity"`
	Price float64 `json:"price"`
}