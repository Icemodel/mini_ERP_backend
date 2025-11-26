package model

import (
	"time"
)

// Role represents the user role values stored in the DB.
type Role int

const (
	RoleAdmin  Role = 1
	RoleStaff  Role = 2
	RoleViewer Role = 3
)

func (r Role) String() string {
	switch r {
	case RoleAdmin:
		return "admin"
	case RoleStaff:
		return "staff"
	case RoleViewer:
		return "viewer"
	default:
		return "unknown"
	}
}

type User struct {
	UserID    string    `gorm:"column:user_id;type:uuid;default:gen_random_uuid();primaryKey" json:"user_id"`
	Username  string    `gorm:"column:username;not null;uniqueIndex" json:"username"`
	FirstName string    `gorm:"column:first_name;not null" json:"first_name"`
	LastName  string    `gorm:"column:last_name;not null" json:"last_name"`
	Password  string    `gorm:"column:password;not null" json:"-"`
	RoleID    Role      `gorm:"column:role_id;not null" json:"role_id"`
	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;autoUpdateTime" json:"updated_at"`
}
