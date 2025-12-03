package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserDataCtx struct {
	UserId uuid.UUID
	Role   string
}

const (
	CONTEXT_USER_DATA_KEY = "__user_data__"
)

// User data
func SetUserDataLocal(c *fiber.Ctx, userData UserDataCtx) {
	c.Locals(CONTEXT_USER_DATA_KEY, userData)
}

func GetUserDataLocal(c *fiber.Ctx) UserDataCtx {
	return c.Locals(CONTEXT_USER_DATA_KEY).(UserDataCtx)
}
