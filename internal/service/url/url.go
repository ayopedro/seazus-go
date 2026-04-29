package url

import (
	"context"
	"errors"

	"github.com/ayopedro/seazus-go/cmd/api/dto"
	appErrors "github.com/ayopedro/seazus-go/internal/common/app_errors"
	"github.com/ayopedro/seazus-go/internal/models"
	"github.com/ayopedro/seazus-go/internal/repository"
	"go.uber.org/zap"
)

type Service interface {
	GetURL(ctx context.Context, id, uID string) (*models.URL, error)
	GetOriginalURL(ctx context.Context, short_url string) (string, error)
	CreateShortURL(ctx context.Context, payload *models.CreateURL, uID string) (string, error)
}

type urlService struct {
	repo   repository.URLRepository
	logger *zap.Logger
}

func NewService(r repository.URLRepository, logger *zap.Logger) Service {
	return &urlService{r, logger}
}

func (us *urlService) GetURL(ctx context.Context, id, uID string) (*models.URL, error) {
	url, err := us.repo.GetOne(ctx, id, uID)

	if err != nil {
		return nil, appErrors.ErrNotFound
	}

	return url, nil
}

func (us *urlService) GetOriginalURL(ctx context.Context, short_url string) (string, error) {
	original, err := us.repo.GetOriginalURL(ctx, short_url)

	if err != nil {
		return "", appErrors.ErrNotFound
	}

	return original, nil
}

func (us *urlService) CreateShortURL(ctx context.Context, payload *models.CreateURL, uID string) (string, error) {
	url := &dto.CreateURLPayload{
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
