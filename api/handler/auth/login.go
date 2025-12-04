package auth

import (
	"log/slog"
	"mini-erp-backend/api/service/auth/command"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

type LoginRequest struct {
	Username string `json:"username" form:"username" query:"username"`
	Password string `json:"password" form:"password" query:"password"`
}

func Login(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := &LoginRequest{}

		if err := c.BodyParser(&req); err != nil {
			logger.Error("invalid request body", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		if req.Username == "" || req.Password == "" {
			logger.Error("invalid login parameters, either username or password is missing")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing email, password or tenantId"})
		}

		request := command.LoginRequest{
			Username: req.Username,
			Password: req.Password,
		}

		response, err := mediatr.Send[*command.LoginRequest, *command.LoginResult](c.Context(), &request)
		if err != nil {
			logger.Error("login command failed", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		// Set refresh token in an HttpOnly secure cookie (do not expose it in JSON)
		if response != nil && response.RefreshToken != "" {
			c.Cookie(&fiber.Cookie{
				Name:     "refresh_token",
				Value:    response.RefreshToken,
				HTTPOnly: true,
				Secure:   false,
				SameSite: "Lax",
				Path:     "/",
				Expires:  time.Unix(response.RefreshTokenExp, 0),
			})

			// prevent returning refresh token in JSON body
			response.RefreshToken = ""
			response.RefreshTokenExp = 0
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}