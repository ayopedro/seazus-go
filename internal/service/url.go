package service

import (
	"context"
	"errors"

	appErrors "github.com/ayopedro/seazus-go/internal/common/app_errors"
	"github.com/ayopedro/seazus-go/internal/models"
	"github.com/ayopedro/seazus-go/internal/repository"
	"go.uber.org/zap"
)

type urlService struct {
	repo   repository.URLRepository
	logger *zap.Logger
}

func NewURLService(r repository.URLRepository, logger *zap.Logger) *urlService {
	return &urlService{r, logger}
}

func (us *urlService) GetURL(ctx context.Context, id, uID string) (*models.URL, error) {
	url, err := us.repo.GetOne(ctx, id, uID)

	if err != nil {
		return nil, appErrors.ErrNotFound
	}

	return url, nil
}

func (us *urlService) CreateShortURL(ctx context.Context, payload *models.CreateURLPayload, uID string) (string, error) {
	url := &models.CreateURLPayload{
		Identifier:  payload.Identifier,
		Url:         payload.Url,
		Description: payload.Description,
	}

	shortUrl, err := us.repo.CreateShortURL(ctx, url, uID)

	if err != nil {
		if errors.Is(err, appErrors.ErrConflict) {
			return "", appErrors.ErrConflict
		}
		return "", appErrors.ErrCreatingShortURL
	}

	return shortUrl, nil
}
