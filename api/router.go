package api

import (
	"log/slog"
	auth_handler "mini-erp-backend/api/handler/auth"
	register_handler "mini-erp-backend/api/handler/register"
	"mini-erp-backend/lib/jwt"
	"mini-erp-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func Register(
	app *fiber.App,
	logger *slog.Logger,
	jwt jwt.Manager,
	mid *middleware.FiberMiddleware,
) {
	v1 := app.Group("/api/v1")

	authGroupApi := v1.Group("/auth")
	{
		authGroupApi.Post("/login", auth_handler.Login(logger))
		authGroupApi.Post("/token/refresh", auth_handler.RefreshLoginToken(logger))
	}
	roleGroupApi := v1.Group("/role")
	{
		roleGroupApi.Use(mid.Authenticated())

		roleGroupApi.Post("/register", mid.RequireMinRole("admin"), register_handler.Register(logger))
	}
}
