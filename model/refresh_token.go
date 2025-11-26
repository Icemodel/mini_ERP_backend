package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	RefreshId  uuid.UUID `gorm:"column:refresh_id;type:uuid;primaryKey" json:"refresh_id"`
	UserNumber uuid.UUID `gorm:"column:user_id;type:uuid;not null;index" json:"user_id"`
	Token      string    `gorm:"column:token;not null;uniqueIndex" json:"token"`
	IssueAt    time.Time `gorm:"column:issue_at;not null" json:"issue_at"`
	ExpireAt   time.Time `gorm:"column:expire_at;not null" json:"expire_at"`
	Revoked    bool      `gorm:"column:revoked;not null" json:"revoked"`
	CreatedAt  time.Time `gorm:"column:created_at;not null;autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null;autoUpdateTime" json:"updated_at"`

	Users User `gorm:"foreignKey:UserId;references:UserNumber;constraint:OnDelete:CASCADE;" json:"-"`
}
