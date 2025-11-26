package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Role represents user role stored as a string in DB and JSON.
type Role string

const (
	Admin  Role = "Admin"
	Staff  Role = "Staff"
	Viewer Role = "Viewer"
)

func (r Role) String() string {
	if r == "" {
		return ""
	}
	return string(r)
}

func (r Role) Value() (driver.Value, error) {
	return string(r), nil
}

func (r *Role) Scan(src interface{}) error {
	if src == nil {
		*r = ""
		return nil
	}
	switch v := src.(type) {
	case string:
		*r = Role(v)
		return nil
	case []byte:
		*r = Role(string(v))
		return nil
	default:
		return fmt.Errorf("cannot scan Role from %T", src)
	}
}

type User struct {
	UserId    string    `gorm:"column:user_id;type:uuid;default:gen_random_uuid();primaryKey" json:"user_id"`
	Username  string    `gorm:"column:username;not null;uniqueIndex" json:"username"`
	FirstName string    `gorm:"column:first_name;not null" json:"first_name"`
	LastName  string    `gorm:"column:last_name;not null" json:"last_name"`
	Password  string    `gorm:"column:password;not null" json:"-"`
	Role      Role      `gorm:"column:role;not null;" json:"role"`
	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;autoUpdateTime" json:"updated_at"`

	AuditLogs []AuditLog `gorm:"foreignKey:UserID;references:UserID" json:"-"`
}
