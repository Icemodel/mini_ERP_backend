package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	AuditLogId uuid.UUID       `gorm:"column:audit_log_id;primaryKey;autoIncrement" json:"audit_log_id"`
	UserId     uuid.UUID       `gorm:"column:user_id; not null" json:"user_id"`
	Action     string          `gorm:"column:action; not null; size:100" json:"action"`
	Detail     json.RawMessage `gorm:"column:detail; not null; type:json" json:"detail"`
	CreatedAt  time.Time       `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`

	User User `gorm:"foreignKey:UserId;references:UserId" json:"-"`
}
