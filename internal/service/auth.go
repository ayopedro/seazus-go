package service

import (
	"context"
	"errors"
	"strings"
	"time"

	utils "github.com/ayopedro/seazus-go/internal/common"
	appErrors "github.com/ayopedro/seazus-go/internal/common/app_errors"
	"github.com/ayopedro/seazus-go/internal/logger"
	"github.com/ayopedro/seazus-go/internal/models"
	"github.com/ayopedro/seazus-go/internal/repository"
	"github.com/google/uuid"
)

type authService struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewAuthService(r repository.UserRepository, jwtSecret string) AuthService {
	return &authService{r, jwtSecret}
}

func (as *authService) CreateUser(ctx context.Context, u *models.CreateUserRequest) error {
	uID, _ := uuid.NewV7()

	hash, err := utils.HashPassword(u.Password)
	if err != nil {
		logger.Error("Error hashing password")
		return appErrors.ErrInternalServerError
	}

	user := &models.User{
		Id:        uID.String(),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     strings.ToLower(u.Email),
		Password:  hash,
	}

	err = as.repo.Create(ctx, user)
	return err
}

func (as *authService) LoginUser(ctx context.Context, p *models.LoginUserRequest) (*models.AuthResponse, error) {
	user, err := as.repo.GetWithEmail(ctx, p.Email)
	if err != nil {
		if errors.Is(err, appErrors.ErrUserNotFound) {
			return nil, appErrors.ErrInvalidCredentials
		}
		return nil, appErrors.ErrInternalServerError
	}

	if user == nil {
		return nil, appErrors.ErrInvalidCredentials
	}

	if err = utils.ComparePasswords(p.Password, user.Password); err != nil {
		return nil, appErrors.ErrInvalidCredentials
	}

	token, _ := utils.GenerateToken(as.jwtSecret, user.Id, 1*time.Hour)

	response := &models.AuthResponse{
		User: models.AuthUser{
			Id:        user.Id,
			FirstName: user.FirstName,
			Email:     user.Email,
		},
		Token: token,
	}

	return response, nil
}
