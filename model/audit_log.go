package model

import (
	"encoding/json"
	"time"
)

type AuditLog struct {
	AuditLogID uint            `gorm:"column:audit_log_id;primaryKey;autoIncrement" json:"audit_log_id"`
	UserID     uint            `gorm:"column:user_id; not null" json:"user_id"`
	Action     string          `gorm:"column:action; not null; size:100" json:"action"`
	Detail     json.RawMessage `gorm:"column:detail; not null; type:json" json:"detail"`
	CreatedAt  time.Time       `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
}
