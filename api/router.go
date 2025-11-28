package api

import (
	"log/slog"
	auth_handler "mini-erp-backend/api/handler/auth"
	"mini-erp-backend/lib/jwt"

	"github.com/gofiber/fiber/v2"
)

func Register(
	app *fiber.App,
	logger *slog.Logger,
	jwt jwt.Manager,
) {
	v1 := app.Group("/api/v1")

	authGroupApi := v1.Group("/auth")
	{
		authGroupApi.Post("/login", auth_handler.Login(logger))
		authGroupApi.Post("/register", auth_handler.Register(logger)) // ใช้ Test เพิ่ม User
	}
}
