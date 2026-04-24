package service

import (
	"context"

	"github.com/ayopedro/seazus-go/internal/models"
)

type UserService interface {
	CreateUser(ctx context.Context, u *models.CreateUserRequest) error
	GetUserProfile(ctx context.Context, uID string) (*models.User, error)
}
