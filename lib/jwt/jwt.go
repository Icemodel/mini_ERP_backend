package jwt

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"mini-erp-backend/config/environment"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type LoginTokenDetail struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
	UserId       uuid.UUID
}

type LoginAccessClaims struct {
	AccessUuid string    `json:"access_uuid"`
	Authorized bool      `json:"authorized"`
	UserId     uuid.UUID `json:"user_id"`
	Role       string    `json:"role"`
	jwt.RegisteredClaims
}

type LoginRefreshClaims struct {
	RefreshUuid string    `json:"refresh_uuid"`
	UserId      uuid.UUID `json:"user_id"`
	Role        string    `json:"role"`
	jwt.RegisteredClaims
}

type loginConfig struct {
	AccessExpMinsLogin  int
	RefreshExpMinsLogin int
	AccessSecret        string
	RefreshSecret       string
}

type manager struct {
	loginConfig loginConfig
	logger      *slog.Logger
}

type Manager interface {
	GenerateLoginToken(userId uuid.UUID, role string) (*LoginTokenDetail, error)
	GetAccessToken(tokenStr string) (*LoginAccessClaims, error)
	GetRefreshToken(tokenStr string) (*LoginRefreshClaims, error)
}

func New(logger *slog.Logger) Manager {
	// fallback: create a default logger when none provided
	if logger == nil {
		handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
		logger = slog.New(handler)
	}

	return &manager{
		loginConfig: loginConfig{
			AccessExpMinsLogin:  environment.GetInt(environment.AccessTokenExpMins),
			RefreshExpMinsLogin: environment.GetInt(environment.RefreshTokenExpMins),
			AccessSecret:        environment.GetString(environment.AccessTokenSecret),
			RefreshSecret:       environment.GetString(environment.RefreshTokenSecret),
		},
		logger: logger,
	}
}

func IsTokenExpired(err error) bool {
	tokenExpiredErr := fmt.Sprintf("%s: %s", jwt.ErrTokenInvalidClaims, jwt.ErrTokenExpired)
	return tokenExpiredErr == err.Error()
}

func (m *manager) createJwt(secret string, claims jwt.Claims) (token string, err error) {
	token, err = jwt.
		NewWithClaims(jwt.SigningMethodHS512, claims).
		SignedString([]byte(secret))
	if err != nil {
		if m.logger != nil {
			m.logger.Error("jwt sign failed", "error", err.Error())
		}
		return token, err
	}

	return token, nil
}

func (m *manager) validateToken(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err := fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			if m.logger != nil {
				m.logger.Error("unexpected signing method", "alg", token.Header["alg"])
			}
			return nil, err
		}

		return []byte(secret), nil
	}
}

func (m *manager) GenerateLoginToken(userId uuid.UUID, role string) (*LoginTokenDetail, error) {
	if userId == uuid.Nil {
		errMsg := "UserID is empty"
		if m.logger != nil {
			m.logger.Error(errMsg)
		}
		return nil, errors.New(errMsg)
	}

	var err error
	token := &LoginTokenDetail{}
	token.UserId = userId

	accessSecret, refreshSecret := m.loginConfig.AccessSecret, m.loginConfig.RefreshSecret
	if accessSecret == "" || refreshSecret == "" {
		errMsg := "token secret from environment is empty"
		m.logger.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	accessExp := time.Now().Add(time.Duration(m.loginConfig.AccessExpMinsLogin) * time.Minute)
	refreshExp := time.Now().Add(time.Duration(m.loginConfig.RefreshExpMinsLogin) * time.Minute)

	// set LoginAccessClaims struct for access token
	accessUUID := uuid.New().String()
	authorized := true
	loginAccessClaims := LoginAccessClaims{
		AccessUuid: accessUUID,
		Authorized: authorized,
		UserId:     userId,
		Role:       role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(accessExp),
		},
	}

	// create access token
	var accessToken string
	accessToken, err = m.createJwt(accessSecret, loginAccessClaims)
	if err != nil {
		m.logger.Error(err.Error())
		return nil, err
	}

	// asign AccessToken token and AtExpires
	token.AccessToken = accessToken
	token.AtExpires = accessExp.Unix()
	token.AccessUuid = accessUUID

	// set LoginRefreshClaims struct for refresh token
	refreshUUID := uuid.New().String()
	loginRefreshClaims := LoginRefreshClaims{
		RefreshUuid: refreshUUID,
		UserId:      userId,
		Role:        role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// create refresh token
	var refreshToken string
	refreshToken, err = m.createJwt(refreshSecret, loginRefreshClaims)
	if err != nil {
		m.logger.Error(err.Error())
		return nil, err
	}

	// asign RefreshToken token and AtExpires
	token.RefreshToken = refreshToken
	token.RtExpires = refreshExp.Unix()
	token.RefreshUuid = refreshUUID

	return token, nil
}

func (m manager) GetAccessToken(tokenStr string) (*LoginAccessClaims, error) {
	secret := environment.GetString(environment.AccessTokenSecret)
	if secret == "" {
		errMsg := "access token secret from environment is empty"
		m.logger.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	token, err := jwt.ParseWithClaims(tokenStr, &LoginAccessClaims{}, m.validateToken(secret))
	if err != nil {
		m.logger.Error("can not parse token", "error", err.Error())
		return nil, err
	}

	if claims, ok := token.Claims.(*LoginAccessClaims); ok {
		return claims, nil
	} else {
		err = errors.New("unknown claims type, cannot proceed")
		m.logger.Error("can not get claims", "error", err.Error())
		return nil, err
	}
}

func (m manager) GetRefreshToken(tokenStr string) (*LoginRefreshClaims, error) {
	secret := environment.GetString(environment.RefreshTokenSecret)
	if secret == "" {
		errMsg := "refresh token secret from environment is empty"
		m.logger.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	token, err := jwt.ParseWithClaims(tokenStr, &LoginRefreshClaims{}, m.validateToken(secret))
	if err != nil {
		m.logger.Error("can not parse token", "error", err.Error())
		return nil, err
	}

	if claims, ok := token.Claims.(*LoginRefreshClaims); ok {
		return claims, nil
	} else {
		err = errors.New("unknown claims type, cannot proceed")
		m.logger.Error("can not get claims", "error", err.Error())
		return nil, err
	}
}
