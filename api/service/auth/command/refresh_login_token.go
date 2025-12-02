package command

import (
	"context"
	"fmt"
	"log/slog"
	"mini-erp-backend/lib/jwt"
	"mini-erp-backend/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshLoginTokenRequest struct {
	AccessToken  string `json:"access_token" form:"access_token" query:"access_token"`
	RefreshToken string `json:"refresh_token" form:"refresh_token" query:"refresh_token"`
}

type RefreshLoginTokenResult struct {
	UserId          uuid.UUID `json:"user_id"`
	AccessToken     string    `json:"access_token"`
	AccessTokenExp  int64     `json:"access_token_exp"`
	RefreshToken    string    `json:"refresh_token"`
	RefreshTokenExp int64     `json:"refresh_token_exp"`
}

type RefreshLoginToken struct {
	domainDb   *gorm.DB
	logger     *slog.Logger
	jwtManager jwt.Manager
	userRepo   repository.User
}

func NewRefreshLoginToken(
	domainDb *gorm.DB,
	logger *slog.Logger,
	jwtManager jwt.Manager,
	userRepo repository.User,
) *RefreshLoginToken {
	return &RefreshLoginToken{
		domainDb:   domainDb,
		logger:     logger,
		jwtManager: jwtManager,
		userRepo:   userRepo,
	}
}

func (r *RefreshLoginToken) Handle(ctx context.Context, request *RefreshLoginTokenRequest) (*RefreshLoginTokenResult, error) {
	var result *RefreshLoginTokenResult

	// Verify current access token exists in DB (user session still valid)
	user, err := r.userRepo.SearchByConditions(r.domainDb, map[string]interface{}{"token": request.AccessToken})
	if err != nil || user == nil {
		if r.logger != nil {
			r.logger.Error("access token not found or invalid", "error", err)
		}
		return nil, fmt.Errorf("invalid access token")
	}

	claims, err := r.jwtManager.ExtractRefreshToken(request.RefreshToken)
	if err != nil {
		if r.logger != nil {
			r.logger.Error("failed to extract refresh token", "error", err)
		}
		if jwt.IsTokenExpired(err) {
			return nil, fmt.Errorf("refresh token expired")
		}
		return nil, fmt.Errorf("invalid refresh token")
	}

	token, err := r.jwtManager.GenerateLoginToken(
		claims.UserId,
		claims.Role,
	)
	if err != nil {
		if r.logger != nil {
			r.logger.Error("failed to generate login token", "error", err)
		}
		return nil, fmt.Errorf("failed to generate login token")
	}

	err = r.userRepo.UpdateTokenByUserId(r.domainDb, claims.UserId, &token.AccessToken)
	if err != nil {
		if r.logger != nil {
			r.logger.Error("failed to update access token in database", "error", err)
		}
		return nil, fmt.Errorf("failed to persist access token")
	}

	result = &RefreshLoginTokenResult{
		UserId:          claims.UserId,
		AccessToken:     token.AccessToken,
		AccessTokenExp:  token.AtExpires,
		RefreshToken:    token.RefreshToken,
		RefreshTokenExp: token.RtExpires,
	}

	return result, nil
}
