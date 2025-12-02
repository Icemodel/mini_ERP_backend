package repository

import (
	"log/slog"
	"mini-erp-backend/model"

	"gorm.io/gorm"
)

type UserRegister interface {
	Create(db *gorm.DB, user model.User) error
}

type userRegister struct {
	logger *slog.Logger
}

func NewUserRegister(logger *slog.Logger) UserRegister {
	return &userRegister{logger: logger}
}

func (r *userRegister) Create(
	db *gorm.DB,
	user model.User,
) error {
	if err := db.Create(&user).Error; err != nil {
		if r.logger != nil {
			r.logger.Error("create user failed", "error", err)
		}
		return err
	}
	return nil
}
