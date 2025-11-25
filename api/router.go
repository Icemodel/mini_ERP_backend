package api

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, logger *slog.Logger) {
	// v1 := app.Group("/api/v1")

	// authGroupApi := v1.Group("/auth")
	// {
	// 	authGroupApi.Post("/register", auth_handler.Register(logger))
	// 	authGroupApi.Post("/login", auth_handler.Login(logger))
	// 	authGroupApi.Post("/refresh", auth_handler.Refresh(logger))
	// 	authGroupApi.Post("/logout", middleware.AuthMiddleware(), auth_handler.Logout(logger))
	// }
}
