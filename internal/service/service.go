package service

import (
	"context"

	"github.com/ayopedro/seazus-go/internal/models"
	"github.com/ayopedro/seazus-go/internal/repository"
	"go.uber.org/zap"
)

type Service struct {
	Auth AuthService
	User UserService
	URL  URLService
}

type UserService interface {
	GetUserProfile(ctx context.Context, uID string) (*models.User, error)
	GetUserURLs(ctx context.Context, uID string) ([]models.URL, error)
}

type AuthService interface {
	CreateUser(ctx context.Context, u *models.CreateUserRequest) error
	LoginUser(ctx context.Context, p *models.LoginUserRequest) (*models.AuthResponse, error)
}

type URLService interface {
	GetURL(ctx context.Context, id, uID string) (*models.URL, error)
	GetOriginalURL(ctx context.Context, short_url string) (string, error)
	CreateShortURL(ctx context.Context, payload *models.CreateURLPayload, uID string) (string, error)
}

func NewService(r *repository.Repository, logger *zap.Logger, jwtSecret string) *Service {
	return &Service{
		Auth: NewAuthService(r.User, jwtSecret, logger),
		User: NewUserService(r.User, r.URL, logger),
		URL:  NewURLService(r.URL, logger),
	}
}
