package command

import (
	"context"
	"errors"
	"log/slog"
	"mini-erp-backend/lib/jwt"
	"mini-erp-backend/model"
	"mini-erp-backend/repository"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Username  string `json:"username" form:"username" query:"username"`
	FirstName string `json:"first_name" form:"first_name" query:"first_name"`
	LastName  string `json:"last_name" form:"last_name" query:"last_name"`
	Password  string `json:"password" form:"password" query:"password"`
	Role      string `json:"role" form:"role" query:"role"`
}

type RegisterResult struct {
	UserId string `json:"user_id"`
}

type UserRegister struct {
	domainDb   *gorm.DB
	logger     *slog.Logger
	jwtManager jwt.Manager
	regisRepo  repository.UserRegister
}

func NewUserRegister(
	domainDb *gorm.DB,
	logger *slog.Logger,
	jwtManager jwt.Manager,
	regisRepo repository.UserRegister,
) *UserRegister {
	return &UserRegister{
		domainDb:   domainDb,
		logger:     logger,
		jwtManager: jwtManager,
		regisRepo:  regisRepo,
	}
}

func (r *UserRegister) Handle(ctx context.Context, request *RegisterRequest) (*RegisterResult, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}

	// validate required fields
	if strings.TrimSpace(request.Username) == "" ||
		strings.TrimSpace(request.Password) == "" ||
		strings.TrimSpace(request.FirstName) == "" ||
		strings.TrimSpace(request.LastName) == "" ||
		strings.TrimSpace(request.Role) == "" {
		return nil, errors.New("missing required fields")
	}

	// normalize inputs
	username := strings.ToLower(strings.TrimSpace(request.Username))
	firstName := strings.ToLower(strings.TrimSpace(request.FirstName))
	lastName := strings.ToLower(strings.TrimSpace(request.LastName))
	roleStr := strings.ToLower(strings.TrimSpace(request.Role))

	// validate role enum
	if roleStr != string(model.RoleAdmin) &&
		roleStr != string(model.RoleStaff) &&
		roleStr != string(model.RoleViewer) {
		return nil, errors.New("invalid role, allowed: admin, staff, viewer")
	}
	role := model.Role(roleStr)

	// check duplicate username
	var existing model.User
	err := r.domainDb.
		Table("users").
		Where("username = ?", username).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("username already exists")
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// OK to create
	} else {
		// unexpected DB error
		return nil, err
	}

	// hash password
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		if r.logger != nil {
			r.logger.Error("password hash failed", "error", err)
		}
		return nil, err
	}

	user := model.User{
		UserId:    uuid.New(),
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Password:  string(hashedPw),
		Role:      role,
	}

	if err := r.regisRepo.Create(r.domainDb, user); err != nil {
		if r.logger != nil {
			r.logger.Error("create user failed", "error", err)
		}
		return nil, err
	}

	return &RegisterResult{UserId: user.UserId.String()}, nil
}
