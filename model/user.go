package model

import (
	"time"

	"github.com/google/uuid"
)

// Role represents user role stored as a string in DB and JSON.

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleStaff  Role = "staff"
	RoleViewer Role = "viewer"
)

type User struct {
	UserId    uuid.UUID `gorm:"column:user_id;type:uuid;primaryKey" json:"user_id"`
	Username  string    `gorm:"column:username;not null;uniqueIndex" json:"username"`
	FirstName string    `gorm:"column:first_name;not null" json:"first_name"`
	LastName  string    `gorm:"column:last_name;not null" json:"last_name"`
	Password  string    `gorm:"column:password;not null" json:"-"`
	Role      Role      `gorm:"column:role; not null;" json:"role"`
	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;autoUpdateTime" json:"updated_at"`
	Token     *string   `gorm:"column:token;" json:"-"`

	AuditLogs []AuditLog `gorm:"foreignKey:UserId;references:UserId" json:"-"`
}
