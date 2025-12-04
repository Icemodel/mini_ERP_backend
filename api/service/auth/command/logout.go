package command

import (
	"context"
	"fmt"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/lib/jwt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LogoutRequest struct {
}

type LogoutResult struct {
	RevokedAt time.Time `json:"revoked_at"`
}

type Logout struct {
	domainDb   *gorm.DB
	logger     *slog.Logger
	jwtManager jwt.Manager
}

type logoutUserCtxKey string

const UserIdContextKey logoutUserCtxKey = "user_id"

func NewLogout(domainDb *gorm.DB, logger *slog.Logger, jwtManager jwt.Manager) *Logout {
	return &Logout{
		domainDb:   domainDb,
		logger:     logger,
		jwtManager: jwtManager,
	}
}

// Handle reads refresh token from context (RefreshTokenContextKey) or user id (UserIdContextKey)
// and revokes sessions for that user.
func (r *Logout) Handle(ctx context.Context, req *LogoutRequest) (*LogoutResult, error) {
	// try refresh token first
	var refreshToken string
	if v := ctx.Value(RefreshTokenContextKey); v != nil {
		if s, ok := v.(string); ok && s != "" {
			refreshToken = s
		}
	}

	sessionRepo := repository.NewUserSession(r.logger)

	if refreshToken != "" {
		// extract claims to get user id
		claims, err := r.jwtManager.ExtractRefreshToken(refreshToken)
		if err != nil {
			if r.logger != nil {
				r.logger.Error("invalid refresh token during logout", "error", err)
			}
			return nil, fmt.Errorf("invalid refresh token")
		}

		// revoke all sessions for this user (single-session semantics)
		if _, err := sessionRepo.RevokeAllByUserID(r.domainDb, claims.UserId); err != nil {
			if r.logger != nil {
				r.logger.Error("failed to revoke sessions during logout", "error", err)
			}
			return nil, fmt.Errorf("failed to revoke sessions")
		}

		return &LogoutResult{RevokedAt: time.Now()}, nil
	}

	// fallback: check user id in context
	if v := ctx.Value(UserIdContextKey); v != nil {
		switch t := v.(type) {
		case uuid.UUID:
			if t != uuid.Nil {
				if _, err := sessionRepo.RevokeAllByUserID(r.domainDb, t); err != nil {
					if r.logger != nil {
						r.logger.Error("failed to revoke sessions by user id", "error", err)
					}
					return nil, fmt.Errorf("failed to revoke sessions")
				}
				return &LogoutResult{RevokedAt: time.Now()}, nil
			}
		case string:
			if t != "" {
				parsed, err := uuid.Parse(t)
				if err == nil && parsed != uuid.Nil {
					if _, err := sessionRepo.RevokeAllByUserID(r.domainDb, parsed); err != nil {
						if r.logger != nil {
							r.logger.Error("failed to revoke sessions by user id", "error", err)
						}
						return nil, fmt.Errorf("failed to revoke sessions")
					}
					return &LogoutResult{RevokedAt: time.Now()}, nil
				}
			}
		}
	}

	// if no token or userid provided
	return nil, fmt.Errorf("missing refresh token or user id")
}