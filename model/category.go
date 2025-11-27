package model

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	CategoryId  uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"category_id"`
	Name        string    `gorm:"not null" json:"name"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null" json:"updated_at"`

	Products []Product `gorm:"foreignKey:CategoryId;references:CategoryId"`
}
