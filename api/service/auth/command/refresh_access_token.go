package command

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/lib/jwt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshAccessTokenRequest struct{}

type RefreshAccessTokenResult struct {
	UserId         uuid.UUID `json:"user_id"`
	SessionId      uuid.UUID `json:"session_id"`
	AccessToken    string    `json:"access_token"`
	AccessTokenExp int64     `json:"access_token_exp"`
}

type RefreshAccessToken struct {
	domainDb   *gorm.DB
	logger     *slog.Logger
	jwtManager jwt.Manager
	userRepo   repository.User
}

type refreshCtxKey string

const RefreshTokenContextKey refreshCtxKey = "refresh_token"

func NewRefreshAccessToken(
	domainDb *gorm.DB,
	logger *slog.Logger,
	jwtManager jwt.Manager,
	userRepo repository.User,
) *RefreshAccessToken {
	return &RefreshAccessToken{
		domainDb:   domainDb,
		logger:     logger,
		jwtManager: jwtManager,
		userRepo:   userRepo,
	}
}

func (r *RefreshAccessToken) Handle(ctx context.Context, request *RefreshAccessTokenRequest) (*RefreshAccessTokenResult, error) {
	var result *RefreshAccessTokenResult
	// Read refresh token from context (set by HTTP handler)
	var refreshToken string
	if v := ctx.Value(RefreshTokenContextKey); v != nil {
		if s, ok := v.(string); ok && s != "" {
			refreshToken = s
		}
	}
	if refreshToken == "" {
		if r.logger != nil {
			r.logger.Error("missing refresh token in context")
		}
		return nil, fmt.Errorf("missing refresh token")
	}

	claims, err := r.jwtManager.ExtractRefreshToken(refreshToken)
	if err != nil {
		if r.logger != nil {
			r.logger.Error("failed to extract refresh token", "error", err)
		}
		if jwt.IsTokenExpired(err) {
			return nil, fmt.Errorf("refresh token expired")
		}
		return nil, fmt.Errorf("invalid refresh token")
	}

	sessionRepo := repository.NewUserSession(r.logger)

	session, err := sessionRepo.GetSessionByRefreshToken(r.domainDb, refreshToken)

	if err != nil {
		if errors.Is(err, repository.ErrSessionRevoked) {
			if r.logger != nil {
				r.logger.Error("refresh token has been revoked", "error", err)
			}
			return nil, fmt.Errorf("refresh token revoked")
		}

		if r.logger != nil {
			r.logger.Error("refresh token not found in session store", "error", err)
		}
		return nil, fmt.Errorf("invalid refresh token")
	}

	if session.UserId != claims.UserId {
		if r.logger != nil {
			r.logger.Error("refresh token user mismatch", "token_user", claims.UserId, "session_user", session.UserId)
		}
		return nil, fmt.Errorf("invalid refresh token")
	}

	accessTokenDetail, err := r.jwtManager.GenerateAccessToken(claims.UserId, claims.Role)
	if err != nil {
		if r.logger != nil {
			r.logger.Error("failed to generate access token", "error", err)
		}
		return nil, fmt.Errorf("failed to generate access token")
	}

	if err := sessionRepo.UpdateAccessToken(r.domainDb, session.SessionId, accessTokenDetail.AccessToken); err != nil {
		if r.logger != nil {
			r.logger.Error("failed to persist new access token", "error", err)
		}
		return nil, fmt.Errorf("failed to persist access token")
	}

	result = &RefreshAccessTokenResult{
		UserId:         claims.UserId,
		SessionId:      session.SessionId,
		AccessToken:    accessTokenDetail.AccessToken,
		AccessTokenExp: accessTokenDetail.AtExpires,
	}

	return result, nil
}
