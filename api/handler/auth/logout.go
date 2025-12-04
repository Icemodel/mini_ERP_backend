package auth

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/service/auth/command"
	"mini-erp-backend/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

// Logout handles session logout. It prefers the HttpOnly refresh cookie; if missing
// and the request is authenticated, it will revoke all sessions for the authenticated user.
func Logout(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Try cookie first
		rt := c.Cookies("refresh_token")
		if rt != "" {
			ctx := context.WithValue(c.Context(), command.RefreshTokenContextKey, rt)
			_, err := mediatr.Send[*command.LogoutRequest, *command.LogoutResult](ctx, &command.LogoutRequest{})
			// clear cookie regardless
			c.Cookie(&fiber.Cookie{Name: "refresh_token", Value: "", Expires: time.Unix(0, 0), Path: "/", HTTPOnly: true})
			if err != nil {
				if logger != nil {
					logger.Error("logout command failed", slog.String("error", err.Error()))
				}
				// keep idempotent: return 204 even on not found/invalid
				return c.SendStatus(fiber.StatusNoContent)
			}
			return c.SendStatus(fiber.StatusNoContent)
		}

		// No cookie: check authenticated user (middleware should set user data)
		// If not authenticated, return 204 idempotent
		ud := c.Locals(utils.CONTEXT_USER_DATA_KEY)
		if ud == nil {
			return c.SendStatus(fiber.StatusNoContent)
		}
		userData, ok := ud.(utils.UserDataCtx)
		if !ok || userData.UserId == uuid.Nil {
			return c.SendStatus(fiber.StatusNoContent)
		}

		// inject user id into context for the Logout command fallback
		ctx := context.WithValue(c.Context(), command.UserIdContextKey, userData.UserId)
		_, err := mediatr.Send[*command.LogoutRequest, *command.LogoutResult](ctx, &command.LogoutRequest{})
		if err != nil {
			if logger != nil {
				logger.Error("logout command failed (by user id)", slog.String("error", err.Error()))
			}
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}
