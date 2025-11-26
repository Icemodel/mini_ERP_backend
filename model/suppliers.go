package model

import (
	"time"
	"github.com/google/uuid"
)

type Supplier struct {
	SupplierId  uuid.UUID `gorm:"type:uuid;primaryKey" json:"supplier_id"`
	Name string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Address string `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}