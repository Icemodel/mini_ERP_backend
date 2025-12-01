package auth

import (
	"fmt"
	"log/slog"
	"mini-erp-backend/api/service/auth/command"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

func RefreshLoginToken(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := &command.RefreshLoginTokenRequest{}

		if err := c.BodyParser(&req); err != nil {
			if logger != nil {
				logger.Error("failed to parse request body", "error", err)
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("invalid request body: %v", err)})
		}

		response, err := mediatr.Send[*command.RefreshLoginTokenRequest, *command.RefreshLoginTokenResult](c.Context(), req)
		if err != nil {
			if logger != nil {
				logger.Error("refresh token command failed", "error", err)
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}
