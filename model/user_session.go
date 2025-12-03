package model

import (
	"time"

	"github.com/google/uuid"
)

type UserSession struct {
	SessionId    uuid.UUID `gorm:"type:uuid;primaryKey" json:"session_id"`
	UserId       uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	AccessToken  string    `gorm:"not null;uniqueIndex" json:"access_token"`
	RefreshToken string    `gorm:"not null;uniqueIndex" json:"refresh_token"`
	Revoked      bool      `gorm:"not null;default:false" json:"revoked"`
	CreatedAt    time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`
}
