package url

import (
	"context"

	"github.com/ayopedro/seazus-go/internal/apperrors"
	"github.com/ayopedro/seazus-go/internal/models"
	url_repository "github.com/ayopedro/seazus-go/internal/repository/url"
	"go.uber.org/zap"
)

type Service interface {
	GetURL(ctx context.Context, id, uID string) (*models.URL, error)
	GetOriginalURL(ctx context.Context, short_url string) (string, error)
	CreateShortURL(ctx context.Context, payload *models.CreateURL, uID string) (string, error)
	UpdateURL(ctx context.Context, id string, payload *models.UpdateURL, uID string) (*models.URL, error)
	DeleteURL(ctx context.Context, id, uID string) error
}

type urlService struct {
	repo   url_repository.Repository
	logger *zap.Logger
}

func NewService(r url_repository.Repository, logger *zap.Logger) Service {
	return &urlService{r, logger}
}

func (us *urlService) GetURL(ctx context.Context, id, uID string) (*models.URL, error) {
	url, err := us.repo.GetOne(ctx, id, uID)

	if err != nil {
		return nil, apperrors.ErrURLNotFound
	}

	return url, nil
}

func (us *urlService) GetOriginalURL(ctx context.Context, short_url string) (string, error) {
	original, err := us.repo.GetOriginalURL(ctx, short_url)

	if err != nil {
		return "", apperrors.ErrURLNotFound
	}

	return original, nil
}

func (us *urlService) CreateShortURL(ctx context.Context, payload *models.CreateURL, uID string) (string, error) {
	shortUrl, err := us.repo.CreateShortURL(ctx, payload, uID)

	if err != nil {
		return "", apperrors.ErrIdentifierTaken
	}

	return shortUrl, nil
}

func (us *urlService) UpdateURL(ctx context.Context, id string, payload *models.UpdateURL, uID string) (*models.URL, error) {

	updatedURL, err := us.repo.UpdateURL(ctx, id, payload, uID)

	if err != nil {
		return nil, err
	}

	return updatedURL, nil
}

func (us *urlService) DeleteURL(ctx context.Context, id, uID string) error {
	err := us.repo.DeleteURL(ctx, id, uID)

	if err != nil {
		return err
	}

	return nil
}
