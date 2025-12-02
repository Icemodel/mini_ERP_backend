package command

import (
	"context"
	"errors"
	"log/slog"
	"mini-erp-backend/lib/jwt"
	"mini-erp-backend/model"
	"mini-erp-backend/repository"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Username string
	Password string
}

type LoginResult struct {
	UserId          uuid.UUID `json:"id"`
	Username        string    `json:"username"`
	Role            string    `json:"role"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	AccessToken     string    `json:"access_token"`
	AccessTokenExp  int64     `json:"access_token_exp"`
	RefreshToken    string    `json:"refresh_token"`
	RefreshTokenExp int64     `json:"refresh_token_exp"`
	SessionId       uuid.UUID `json:"session_id"`
}

type LoginByUsername struct {
	domainDb   *gorm.DB
	logger     *slog.Logger
	jwtManager jwt.Manager
	userRepo   repository.User
}

func NewLoginByUsername(
	domainDb *gorm.DB,
	logger *slog.Logger,
	jwtManager jwt.Manager,
	userRepo repository.User,
) *LoginByUsername {
	return &LoginByUsername{
		domainDb:   domainDb,
		logger:     logger,
		jwtManager: jwtManager,
		userRepo:   userRepo,
	}
}

func (l *LoginByUsername) Handle(ctx context.Context, request *LoginRequest) (*LoginResult, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}

	username := strings.TrimSpace(strings.ToLower(request.Username))
	if username == "" || request.Password == "" {
		return nil, errors.New("username and password are required")
	}

	user, err := l.userRepo.Search(l.domainDb, username)
	if err != nil {
		if l.logger != nil {
			l.logger.Error("user lookup failed", "username", username, "error", err)
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		if l.logger != nil {
			l.logger.Error("invalid credentials", "username", username)
		}
		return nil, err
	}

	roleStr := string(user.Role)
	token, err := l.jwtManager.GenerateLoginToken(user.UserId, roleStr)
	if err != nil {
		if l.logger != nil {
			l.logger.Error("generate token failed", "user", user.UserId.String(), "error", err)
		}
		return nil, err
	}

	session := &model.UserSession{
		SessionId:    uuid.New(),
		UserId:       user.UserId,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Revoked:      false,
	}

	sessRepo := repository.NewUserSession(l.logger)
	if err := sessRepo.Create(l.domainDb, session); err != nil {
		if l.logger != nil {
			l.logger.Error("create user session failed", "user", user.UserId.String(), "error", err)
		}
		return nil, err
	}

	res := &LoginResult{
		UserId:          user.UserId,
		Username:        user.Username,
		Role:            string(user.Role),
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		AccessToken:     token.AccessToken,
		AccessTokenExp:  token.AtExpires,
		RefreshToken:    token.RefreshToken,
		RefreshTokenExp: token.RtExpires,
		SessionId:       session.SessionId,
	}

	return res, nil
}
