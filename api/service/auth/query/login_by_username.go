package query

import (
	"context"
	"errors"
	"log/slog"
	"mini-erp-backend/lib/jwt"
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
	UserId   uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
}

type LoginByUsername struct {
	domainDb   *gorm.DB
	logger     *slog.Logger
	jwtManager jwt.Manager
	userRepo   repository.UserAuthen
}

func NewLoginByUsername(
	domainDb *gorm.DB,
	logger *slog.Logger,
	jwtManager jwt.Manager,
	userRepo repository.UserAuthen,
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

	creds, err := l.userRepo.Search(l.domainDb, username)
	if err != nil {
		if l.logger != nil {
			l.logger.Error("user lookup failed", "username", username, "error", err)
		}
		return nil, err
	}

	//ต้องมาแก้ให้ใช้ bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(creds.Password), []byte(request.Password)); err != nil {
		if l.logger != nil {
			l.logger.Error("invalid credentials", "username", username)
		}
		return nil, err
	}

	// ต้องCopy การให้ Token อีก อันนี้ Place holder
	roleStr := creds.Role.String()
	_, err = l.jwtManager.GenerateLoginToken(creds.UserId, roleStr)
	if err != nil {
		if l.logger != nil {
			l.logger.Error("generate token failed", "user", creds.UserId.String(), "error", err)
		}
		return nil, err
	}

	res := &LoginResult{
		UserId:   creds.UserId,
		Username: creds.Username,
		Role:     roleStr,
	}

	return res, nil
}
