package auth

import (
	"log/slog"
	"mini-erp-backend/api/service/auth/query"
	"mini-erp-backend/lib/jwt"
	"mini-erp-backend/repository"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func NewService(
	domainDb *gorm.DB,
	logger *slog.Logger,
	jwtManager jwt.Manager,
	userRepo repository.UserAuthen,
) {
	LoginService := query.NewLoginByUsername(
		domainDb,
		logger,
		jwtManager,
		userRepo,
	)

	err := mediatr.RegisterRequestHandler(LoginService)
	if err != nil {
		panic(err)
	}
}
