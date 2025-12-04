package register

import (
	"log/slog"
	"mini-erp-backend/api/service/register/command"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

type RegisterRequest struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	Role      string `json:"role"`
}

func Register(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := &RegisterRequest{}

		if err := c.BodyParser(&req); err != nil {
			logger.Error("invalid request body", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		request := command.RegisterRequest{
			Username:  req.Username,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Password:  req.Password,
			Role:      req.Role,
		}

		response, err := mediatr.Send[*command.RegisterRequest, *command.RegisterResult](c.Context(), &request)
		if err != nil {
			logger.Error("login command failed", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}
