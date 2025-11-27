package auth

import (
	"log/slog"
	"mini-erp-backend/api/service/auth/query"

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
			logger.Error("invalid login parameters")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing email, password or tenantId"})
		}

		request := query.LoginRequest{
			Username: req.Username,
			Password: req.Password,
		}

		response, err := mediatr.Send[*query.LoginRequest, *query.LoginResult](c.Context(), &request)
		if err != nil {
			logger.Error("login command failed", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}
