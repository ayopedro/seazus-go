package service

import (
	"context"

	"github.com/ayopedro/seazus-go/internal/models"
)

type UserService interface {
	GetUserProfile(ctx context.Context, uID string) (*models.User, error)
}

type AuthService interface {
	CreateUser(ctx context.Context, u *models.CreateUserRequest) error
	LoginUser(ctx context.Context, p *models.LoginUserRequest) (*models.AuthResponse, error)
}

type URLService interface {
	GetURL(ctx context.Context, id, uID string) (*models.URL, error)
	// CreateShortURL(ctx context.Context, payload *models.CreateURLPayload) (string, error)
}
