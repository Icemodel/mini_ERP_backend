package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
    AuditLogId uuid.UUID       `gorm:"type:uuid;primaryKey" json:"audit_log_id"`
    UserId     uuid.UUID       `gorm:"type:uuid;not null" json:"user_id"`
    Action     string          `gorm:"size:100;not null" json:"action"`
    Detail     json.RawMessage `gorm:"type:json;not null" json:"detail"`
    CreatedAt  time.Time       `gorm:"not null;autoCreateTime" json:"created_at"`

    User User `gorm:"foreignKey:UserId;references:UserId" json:"-"`
}

