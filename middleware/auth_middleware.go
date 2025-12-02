package middleware

import (
	"log/slog"
	"mini-erp-backend/config/environment"
	"mini-erp-backend/lib/jwt"
	"mini-erp-backend/repository"
	"mini-erp-backend/utils"
	"os"
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

type FiberMiddleware struct {
	db         *gorm.DB
	corsSetUp  corsSetUp
	logger     *slog.Logger
	jwtManager jwt.Manager
	userRepo   repository.UserAuthen
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
	userRepo repository.UserAuthen,
) *FiberMiddleware {
	if fiberMiddlewareInstance == nil {
		fiberMiddlewareLock.Lock()
		defer fiberMiddlewareLock.Unlock()
		if fiberMiddlewareInstance == nil {
			fiberMiddlewareInstance = createFiberMiddlewareInstance(
				db,
				logger,
				jwtManager,
				userRepo,
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
	userRepo repository.UserAuthen,
) *FiberMiddleware {
	return getFiberMiddlewareInstance(db, logger, jwtManager, userRepo)
}

// create the fiberMiddlewareInstance and set up it
func createFiberMiddlewareInstance(
	db *gorm.DB,
	logger *slog.Logger,
	jwtManager jwt.Manager,
	userRepo repository.UserAuthen,
) *FiberMiddleware {
	allowCredential, err := strconv.ParseBool(environment.GetString(environment.AllowCredentialKey))
	if err != nil {
		message := "Failed to set CORS config"
		logger.Error(message, err.Error())
		os.Exit(1)
	}

	return &FiberMiddleware{
		corsSetUp: corsSetUp{
			AllowOrigins:     environment.GetString(environment.AllowOriginKey),
			AllowCredentials: allowCredential,
		},
		db:         db,
		logger:     logger,
		jwtManager: jwtManager,
		userRepo:   userRepo,
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
				return c.Status(fiber.StatusBadRequest).JSON(err.Error())
			} else {
				f.logger.Error(err.Error())
				return c.Status(fiber.StatusBadRequest).JSON(err.Error())
			}
		}

		conditions := map[string]interface{}{
			"user_id": claims.UserId,
		}

		userDataDetail, err := f.userRepo.SearchByConditions(f.db, conditions)
		if err != nil {
			f.logger.Error(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(err.Error())
		}

		if userDataDetail.Token != nil {
			if tokenStr != *userDataDetail.Token {
				f.logger.Error("token mismatch")
				return c.Status(fiber.StatusBadRequest).JSON("token mismatch")
			}
		} else {
			f.logger.Error("token is nil")
			return c.Status(fiber.StatusBadRequest).JSON("token is nil")
		}

		userData := utils.UserDataCtx{
			UserId: claims.UserId,
			Role:   claims.Role,
		}
		utils.SetUserDataLocal(c, userData)

		return c.Next()
	}
}
