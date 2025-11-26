package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	AuditLogId uuid.UUID       `gorm:"type:uuid;primaryKey" json:"audit_log_id"`
	UserId     uuid.UUID       `gorm:"type:uuid;not null" json:"user_id"`
	Action     string          `gorm:"not null" json:"action"`
	Detail     json.RawMessage `gorm:"type:json;not null" json:"detail"`
	CreatedAt  time.Time       `gorm:"not null;autoCreateTime" json:"created_at"`

	User User `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
}
