package repository

import (
	"log/slog"
	"mini-erp-backend/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSession interface {
	Create(db *gorm.DB, session *model.UserSession) error
	GetSessionByRefreshToken(db *gorm.DB, refreshToken string) (*model.UserSession, error)
	UpdateAccessToken(db *gorm.DB, sessionId uuid.UUID, accessToken string) error
}

type userSession struct {
	logger *slog.Logger
}

func NewUserSession(logger *slog.Logger) UserSession {
	return &userSession{logger: logger}
}

func (r *userSession) Create(db *gorm.DB, session *model.UserSession) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.UserSession{}).
			Where("user_id = ? AND revoked = ?", session.UserId, false).
			Update("revoked", true).Error; err != nil {
			if r.logger != nil {
				r.logger.Error("failed to revoke old user sessions", "user_id", session.UserId, "error", err)
			}
			return err
		}

		if err := tx.Create(session).Error; err != nil {
			if r.logger != nil {
				r.logger.Error("failed to create user session", "error", err)
			}
			return err
		}
		return nil
	})
	return err
}

func (r *userSession) GetSessionByRefreshToken(db *gorm.DB, refreshToken string) (*model.UserSession, error) {
	var session model.UserSession
	err := db.Model(&model.UserSession{}).
		Select("refresh_token", "revoked", "user_id", "session_id").
		Where("refresh_token = ? AND revoked = ?", refreshToken, false).
		First(&session).Error

	if err != nil {
		if r.logger != nil {
			r.logger.Error("failed to find user session by refresh token", "error", err)
		}
		return nil, err
	}

	return &session, nil
}

func (r *userSession) UpdateAccessToken(db *gorm.DB, sessionId uuid.UUID, accessToken string) error {
	// Only update the access_token. Refresh token rotation is handled explicitly elsewhere.
	if err := db.Model(&model.UserSession{}).
		Where("session_id = ?", sessionId).
		Update("access_token", accessToken).Error; err != nil {
		if r.logger != nil {
			r.logger.Error("failed to update access token for session", "session_id", sessionId, "error", err)
		}
		return err
	}

	return nil
}
