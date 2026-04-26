package service

import (
	"context"
	"net/http"

	"github.com/ayopedro/seazus-go/internal/models"
)

type UserService interface {
	GetUserProfile(ctx context.Context, uID string) (*models.User, error)
	GetUserURLs(ctx context.Context, uID string) ([]models.URL, error)
}

type AuthService interface {
	CreateUser(ctx context.Context, u *models.CreateUserRequest) error
	LoginUser(w http.ResponseWriter, ctx context.Context, p *models.LoginUserRequest) (*models.AuthResponse, error)
}

type URLService interface {
	GetURL(ctx context.Context, id, uID string) (*models.URL, error)
	CreateShortURL(ctx context.Context, payload *models.CreateURLPayload, uID string) (string, error)
}
