package jwt

import (
	"mini-erp-backend/config/environment"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserId    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	TokenType string    `json:"token_type"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userId uuid.UUID, username string) (string, int64, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		UserId:    userId,
		Username:  username,
		TokenType: "Access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(environment.GetString("JWT_SECRET")))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expirationTime.Unix(), nil
}

func GenerateRefreshToken(userId uuid.UUID, username string) (string, int64, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId:    userId,
		Username:  username,
		TokenType: "Refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(environment.GetString("JWT_SECRET")))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expirationTime.Unix(), nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(environment.GetString("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	// Check if token type is Access
	if claims.TokenType != "Access" {
		return nil, jwt.ErrInvalidType
	}

	return claims, nil
}

func ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(environment.GetString("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	// Check if token type is Refresh
	if claims.TokenType != "Refresh" {
		return nil, jwt.ErrInvalidType
	}

	return claims, nil
}
