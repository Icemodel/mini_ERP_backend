package auth

import (
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/api/service/auth/command"
	"mini-erp-backend/lib/jwt"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func NewService(
	domainDb *gorm.DB,
	logger *slog.Logger,
	jwtManager jwt.Manager,
	userRepo repository.User,
) {
	LoginService := command.NewLoginByUsername(
		domainDb,
		logger,
		jwtManager,
		userRepo,
	)
	RefreshLoginTokenService := command.NewRefreshAccessToken(
		domainDb,
		logger,
		jwtManager,
		userRepo,
	)

	err := mediatr.RegisterRequestHandler(LoginService)
	if err != nil {
		panic(err)
	}

	err = mediatr.RegisterRequestHandler(RefreshLoginTokenService)
	if err != nil {
		panic(err)
	}
}
