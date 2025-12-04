package repository

import (
	"errors"
	"log/slog"
	"mini-erp-backend/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrSessionRevoked = errors.New("session revoked")

type UserSession interface {
	Create(db *gorm.DB, session *model.UserSession) error
	GetSessionByRefreshToken(db *gorm.DB, refreshToken string) (*model.UserSession, error)
	UpdateAccessToken(db *gorm.DB, sessionId uuid.UUID, accessToken string) error
	RevokeAllByUserID(db *gorm.DB, userId uuid.UUID) (int64, error)
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
		Where("refresh_token = ?", refreshToken).
		First(&session).Error

	if err != nil {
		if r.logger != nil {
			r.logger.Error("failed to find user session by refresh token", "error", err)
		}
		return nil, err
	}

	if session.Revoked {
		if r.logger != nil {
			r.logger.Warn("refresh token found but revoked", "session_id", session.SessionId)
		}
		return nil, ErrSessionRevoked
	}

	return &session, nil
}

func (r *userSession) UpdateAccessToken(db *gorm.DB, sessionId uuid.UUID, accessToken string) error {
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

func (r *userSession) RevokeAllByUserID(db *gorm.DB, userId uuid.UUID) (int64, error) {
	res := db.Model(&model.UserSession{}).
		Where("user_id = ? AND revoked = ?", userId, false).
		Update("revoked", true)
	if res.Error != nil {
		if r.logger != nil {
			r.logger.Error("failed to revoke all sessions for user", "user_id", userId, "error", res.Error)
		}
		return 0, res.Error
	}
	if r.logger != nil {
		r.logger.Info("revoked sessions for user", "user_id", userId, "rows_affected", res.RowsAffected)
	}
	return res.RowsAffected, nil
}
