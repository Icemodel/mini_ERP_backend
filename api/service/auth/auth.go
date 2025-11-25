package auth

import (
	"log/slog"

	"gorm.io/gorm"
)

func NewService(logger *slog.Logger, db *gorm.DB) {
	// loginService := command.NewLogin(logger, db, userRepo, tokenRepo)

	// err := mediatr.RegisterRequestHandler(loginService)
	// if err != nil {
	// 	panic(err)
	// }
}
