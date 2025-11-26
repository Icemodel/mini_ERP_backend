package api

import (
	"log/slog"
	category_handler "mini-erp-backend/api/handler/category"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, logger *slog.Logger) {
	v1 := app.Group("/api/v1")

	categoryGroupApi := v1.Group("/categories")
	{
		categoryGroupApi.Get("/", category_handler.Categories(logger))
		categoryGroupApi.Get("/:id", category_handler.CategoryById(logger))
		categoryGroupApi.Post("/", category_handler.Create(logger))
		categoryGroupApi.Patch("/", category_handler.Update(logger))
		categoryGroupApi.Delete("/:id", category_handler.DeleteById(logger))
	}

	// authGroupApi := v1.Group("/auth")
	// {
	// 	authGroupApi.Post("/register", auth_handler.Register(logger))
	// 	authGroupApi.Post("/login", auth_handler.Login(logger))
	// 	authGroupApi.Post("/refresh", auth_handler.Refresh(logger))
	// 	authGroupApi.Post("/logout", middleware.AuthMiddleware(), auth_handler.Logout(logger))
	// }
}
