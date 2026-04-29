package repository

import (
	"context"
	"database/sql"

	"github.com/ayopedro/seazus-go/cmd/api/dto"
	"github.com/ayopedro/seazus-go/internal/models"
	"go.uber.org/zap"
)

type Repository struct {
	URL  URLRepository
	User UserRepository
}

type URLRepository interface {
	GetOne(ctx context.Context, id, uID string) (*models.URL, error)
	GetUserURLs(ctx context.Context, uID string) ([]models.URL, error)
	GetOriginalURL(ctx context.Context, short_url string) (string, error)
	CreateShortURL(ctx context.Context, payload *dto.CreateURLPayload, uID string) (string, error)
}

type UserRepository interface {
	Create(ctx context.Context, u *models.User) error
	Get(ctx context.Context, uId string) (*models.User, error)
	GetWithEmail(ctx context.Context, email string) (*models.User, error)
}

func NewRepository(db *sql.DB, logger *zap.Logger) *Repository {
	return &Repository{
		URL:  NewURLRepository(db, logger),
		User: NewUserRepository(db, logger),
	}
}
