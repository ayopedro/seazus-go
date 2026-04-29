package auth

import (
	"context"
	"strings"
	"time"

	"github.com/ayopedro/seazus-go/internal/models"
	"github.com/ayopedro/seazus-go/internal/repository/user"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	CreateUser(ctx context.Context, u *models.CreateUser) error
	LoginUser(ctx context.Context, p *models.LoginUser) (*models.AuthResponse, error)
}

type authService struct {
	repo      user.Repository
	jwtSecret string
	logger    *zap.Logger
}

func NewService(r user.Repository, jwtSecret string, logger *zap.Logger) Service {
	return &authService{r, jwtSecret, logger}
}

func (as *authService) CreateUser(ctx context.Context, u *models.CreateUser) error {
	uID, _ := uuid.NewV7()

	hash, err := hashPassword(u.Password)
	if err != nil {
		if as.logger != nil {
			as.logger.Error("Error hashing password", zap.Error(err))
		}
		return err
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

func (as *authService) LoginUser(ctx context.Context, p *models.LoginUser) (*models.AuthResponse, error) {
	user, err := as.repo.GetWithEmail(ctx, p.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, err
	}

	if err = comparePasswords(p.Password, user.Password); err != nil {
		return nil, err
	}

	token, _ := generateToken(as.jwtSecret, user.Id, 1*time.Hour)

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
