package middleware

import (
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/config/environment"
	"mini-erp-backend/lib/jwt"
	"mini-erp-backend/utils"
	"os"
	"strconv"
	"strings"

	"sync"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

type FiberMiddleware struct {
	db          *gorm.DB
	corsSetUp   corsSetUp
	logger      *slog.Logger
	jwtManager  jwt.Manager
	userRepo    repository.User
	sessionRepo repository.UserSession
}

type corsSetUp struct {
	AllowOrigins     string
	AllowCredentials bool
}

// a mutex for synchronizing access to the fiberMiddlewareInstance variable
var fiberMiddlewareLock = &sync.Mutex{}

// a singleton instance of the FiberMiddleware struct
var fiberMiddlewareInstance *FiberMiddleware

// return the singleton instance of the FiberMiddleware
func getFiberMiddlewareInstance(
	db *gorm.DB,
	logger *slog.Logger,
	jwtManager jwt.Manager,
	userAuthenRepo repository.User,
	sessionRepo repository.UserSession,
) *FiberMiddleware {
	if fiberMiddlewareInstance == nil {
		fiberMiddlewareLock.Lock()
		defer fiberMiddlewareLock.Unlock()
		if fiberMiddlewareInstance == nil {
			fiberMiddlewareInstance = createFiberMiddlewareInstance(
				db,
				logger,
				jwtManager,
				userAuthenRepo,
				sessionRepo,
			)
		}
	}

	return fiberMiddlewareInstance
}

// new the fiberMiddlewareInstance and return it out
func NewFiberMiddleware(
	db *gorm.DB,
	logger *slog.Logger,
	jwtManager jwt.Manager,
	userAuthenRepo repository.User,
	sessionRepo repository.UserSession,
) *FiberMiddleware {
	return getFiberMiddlewareInstance(db, logger, jwtManager, userAuthenRepo, sessionRepo)
}

// create the fiberMiddlewareInstance and set up it
func createFiberMiddlewareInstance(
	db *gorm.DB,
	logger *slog.Logger,
	jwtManager jwt.Manager,
	userRepo repository.User,
	sessionRepo repository.UserSession,
) *FiberMiddleware {
	allowCredential, err := strconv.ParseBool(environment.GetString(environment.AllowCredentialKey))
	if err != nil {
		message := "Failed to set CORS config"
		logger.Error(message, "error", err)
		os.Exit(1)
	}

	return &FiberMiddleware{
		corsSetUp: corsSetUp{
			AllowOrigins:     environment.GetString(environment.AllowOriginKey),
			AllowCredentials: allowCredential,
		},
		db:          db,
		logger:      logger,
		jwtManager:  jwtManager,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

// allows servers to specify who can access its resources and what resources can access
func (f *FiberMiddleware) CORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     f.corsSetUp.AllowOrigins,
		AllowCredentials: f.corsSetUp.AllowCredentials,
	})
}

func (f *FiberMiddleware) Authenticated() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr, err := f.jwtManager.GetAccessTokenFromContext(c)
		if err != nil {
			f.logger.Error(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(err.Error())
		}

		claims, err := f.jwtManager.ExtractAccessToken(tokenStr)
		if err != nil {
			if jwt.IsTokenExpired(err) {
				f.logger.Error(err.Error())
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "access token expired"})
			}

			f.logger.Error(err.Error())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid access token"})
		}

		userData := utils.UserDataCtx{
			UserId: claims.UserId,
			Role:   strings.ToLower(claims.Role),
		}
		utils.SetUserDataLocal(c, userData)

		return c.Next()
	}
}

var roleLevel = map[string]int{
	"viewer": 1,
	"staff":  2,
	"admin":  3,
}

func (f *FiberMiddleware) RequireMinRole(minRole string) fiber.Handler {
	min := 999
	if lvl, ok := roleLevel[strings.ToLower(minRole)]; ok {
		min = lvl
	}

	return func(c *fiber.Ctx) error {
		v := c.Locals(utils.CONTEXT_USER_DATA_KEY)
		if v == nil {
			f.logger.Error("authorization: user data missing in context")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		ud, ok := v.(utils.UserDataCtx)
		if !ok {
			f.logger.Error("authorization: invalid user data type in context")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		role := strings.ToLower(strings.TrimSpace(ud.Role))
		if role == "" || ud.UserId == uuid.Nil {
			f.logger.Error("authorization: role or user id missing in user data")
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "role missing or invalid"})
		}

		if lvl, ok := roleLevel[role]; ok && lvl >= min {
			return c.Next()
		}

		f.logger.Error("authorization: forbidden", "required_min_role", minRole, "user_role", role)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden role for this action"})
	}
}

func (f *FiberMiddleware) RequireRole(roles ...string) fiber.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		rn := strings.ToLower(strings.TrimSpace(r))
		if rn != "" {
			allowed[rn] = struct{}{}
		}
	}

	return func(c *fiber.Ctx) error {
		v := c.Locals(utils.CONTEXT_USER_DATA_KEY)
		if v == nil {
			f.logger.Error("authorization: user data missing in context")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		ud, ok := v.(utils.UserDataCtx)
		if !ok {
			f.logger.Error("authorization: invalid user data type in context")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		role := strings.ToLower(strings.TrimSpace(ud.Role))
		if role == "" || ud.UserId == uuid.Nil {
			f.logger.Error("authorization: role or user id missing in user data")
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "role missing or invalid"})
		}

		if _, ok := allowed[role]; ok {
			return c.Next()
		}

		f.logger.Error("authorization: forbidden", "required_roles", roles, "user_role", role)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden role for this action"})
	}
}
