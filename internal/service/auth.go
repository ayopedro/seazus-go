package service

import (
	"context"
	"strings"

	"github.com/ayopedro/seazus-go/internal/logger"
	"github.com/ayopedro/seazus-go/internal/models"
	"github.com/ayopedro/seazus-go/internal/repository"
	"github.com/ayopedro/seazus-go/internal/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type authService struct {
	repo *repository.UserRepository
}

func NewAuthService(r *repository.UserRepository) *authService {
	return &authService{r}
}

func (as *authService) CreateUser(ctx context.Context, u *models.CreateUserRequest) error {
	uID, _ := uuid.NewV7()

	hash, err := utils.HashPassword(u.Password)
	if err != nil {
		logger.Error("Error hashing password")
		return models.ErrInternalServerError
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

func (as *authService) LoginUser(ctx context.Context, p *models.LoginUserRequest) (*models.AuthUser, error) {
	user, err := as.repo.GetWithEmail(ctx, p.Email)
	if err != nil {
		logger.Error("Error getting user:", zap.String("login user", err.Error()))
		return nil, models.ErrInternalServerError
	}

	if user == nil {
		return nil, models.ErrInvalidCredentials
	}

	if err = utils.ComparePasswords(p.Password, user.Password); err != nil {
		return nil, models.ErrInvalidCredentials
	}

	response := &models.AuthUser{
		Id:        user.Id,
		FirstName: user.FirstName,
		Email:     user.Email,
		Token:     "random-access-token",
	}

	return response, nil
}
