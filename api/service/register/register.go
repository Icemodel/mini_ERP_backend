package register

import (
	"log/slog"
	"mini-erp-backend/api/service/register/command"
	"mini-erp-backend/lib/jwt"
	"mini-erp-backend/repository"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func NewService(
	domainDb *gorm.DB,
	logger *slog.Logger,
	jwtManager jwt.Manager,
	regisRepo repository.User,
) {
	RegisterService := command.NewUserRegister(
		domainDb,
		logger,
		jwtManager,
		regisRepo,
	)

	err := mediatr.RegisterRequestHandler(RegisterService)
	if err != nil {
		panic(err)
	}
}
