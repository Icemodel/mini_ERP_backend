package repository

import (
	"log/slog"
	"strings"

	"mini-erp-backend/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserAuthen interface {
	Search(db *gorm.DB, username string) (*model.User, error)
	SearchByConditions(db *gorm.DB, conditions map[string]interface{}) (*model.User, error)
	UpdateTokenByUserId(db *gorm.DB, userId uuid.UUID, token *string) error
	SearchUserByToken(db *gorm.DB, token string) (*model.User, error)
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

func (r *userAuthen) SearchByConditions(db *gorm.DB, conditions map[string]interface{}) (*model.User, error) {
	var user model.User

	if err := db.
		Table("users").
		Select("user_id", "username", "first_name", "last_name", "password", "role", "token", "created_at", "updated_at").
		Where(conditions).
		First(&user).Error; err != nil {
		if r.logger != nil {
			r.logger.Error("query user by conditions failed", "conditions", conditions, "error", err)
		}
		return nil, err
	}

	return &user, nil
}

func (r userAuthen) UpdateTokenByUserId(db *gorm.DB, userId uuid.UUID, token *string) error {
	if err := db.Model(&model.User{}).
		Where("user_id = ?", userId).
		Update("token", token).
		Error; err != nil {
		if r.logger != nil {
			r.logger.Error("can not update token by user id", "user_id", userId, "error", err)
		}
		return err
	}
	return nil
}

func (r *userAuthen) SearchUserByToken(db *gorm.DB, token string) (*model.User, error) {
	var user model.User

	if err := db.
		Where("token = ?", token).
		First(&user).
		Error; err != nil {
		if r.logger != nil {
			r.logger.Error("can not find user by token", "token", token, "error", err)
		}
		return nil, err
	}

	return &user, nil
}
