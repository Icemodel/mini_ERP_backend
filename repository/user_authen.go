package repository

import (
	"log/slog"
	"strings"

	"mini-erp-backend/model"

	"gorm.io/gorm"
)

type UserAuthen interface {
	Search(db *gorm.DB, username string) (*model.User, error)
}

type userAuthen struct {
	logger *slog.Logger
}

func NewUserAuthen(logger *slog.Logger) UserAuthen {
	return &userAuthen{logger: logger}
}

func (r *userAuthen) Search(db *gorm.DB, username string) (*model.User, error) {
	var user model.User
	username = strings.TrimSpace(strings.ToLower(username))
	if err := db.
		Table("users").
		Select("user_id", "username", "first_name", "last_name", "password", "role", "created_at", "updated_at").
		Where("username = ?", username).
		First(&user).Error; err != nil {
		if r.logger != nil {
			r.logger.Error("query user credentials failed", "username", username, "error", err)
		}
		return nil, err
	}
	return &user, nil
}
