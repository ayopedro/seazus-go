package repository

import (
	"context"

	"github.com/ayopedro/seazus-go/internal/models"
)

type URLRepository interface {
	GetOne(ctx context.Context, id, uID string) (*models.URL, error)
}

type UserRepository interface {
	Create(ctx context.Context, u *models.User) error
	Get(ctx context.Context, uId string) (*models.User, error)
	GetWithEmail(ctx context.Context, email string) (*models.User, error)
}
