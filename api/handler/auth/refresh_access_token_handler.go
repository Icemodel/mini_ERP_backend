package auth

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/service/auth/command"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

func RefreshAccessToken(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Read refresh token only from HttpOnly cookie (frontend must call withCredentials)
		rt := c.Cookies("refresh_token")
		if rt == "" {
			if logger != nil {
				logger.Error("missing refresh token cookie")
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing refresh token"})
		}

		// Pass refresh token only via context (cookie). Command will read it from context.
		req := &command.RefreshAccessTokenRequest{}
		ctx := context.WithValue(c.Context(), command.RefreshTokenContextKey, rt)

		response, err := mediatr.Send[*command.RefreshAccessTokenRequest, *command.RefreshAccessTokenResult](ctx, req)
		if err != nil {
			if logger != nil {
				logger.Error("refresh token command failed", "error", err)
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}
