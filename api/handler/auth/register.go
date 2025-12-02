package auth

import (
	"log/slog"
	"mini-erp-backend/config/database"
	"mini-erp-backend/config/environment"
	"mini-erp-backend/model"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	Role      string `json:"role"`
}

// Register creates a new user. This is a minimal example and does not
// perform password hashing or advanced validation.
func Register(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := &RegisterRequest{}

		if err := c.BodyParser(req); err != nil {
			if logger != nil {
				logger.Error("invalid register body", "error", err)
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		if req.Username == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing required fields"})
		}

		db := database.Connect(environment.GetString("DSN_DATABASE"))

		// hash password before storing
		hashedPw, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			if logger != nil {
				logger.Error("password hash failed", "error", err)
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to process password"})
		}

		user := model.User{
			UserId:    uuid.New(),
			Username:  strings.ToLower(strings.TrimSpace(req.Username)),
			FirstName: strings.ToLower(strings.TrimSpace(req.FirstName)),
			LastName:  strings.ToLower(strings.TrimSpace(req.LastName)),
			Password:  string(hashedPw),
			Role:      string(req.Role),
		}

		if err := db.Create(&user).Error; err != nil {
			if logger != nil {
				logger.Error("create user failed", "error", err)
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create user"})
		}

		// Password field has `json:"-"` in the model so it won't be returned
		return c.Status(fiber.StatusCreated).JSON(user)
	}
}
