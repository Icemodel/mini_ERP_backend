package model

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	CategoryId  uuid.UUID `gorm:"type:uuid;primaryKey" json:"category_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Products []Product `gorm:"foreignKey:CategoryId;references:CategoryId"`
}
