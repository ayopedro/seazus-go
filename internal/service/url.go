package service

import (
	"context"

	"github.com/ayopedro/seazus-go/internal/models"
	"github.com/ayopedro/seazus-go/internal/repository"
)

type urlService struct {
	repo *repository.URLRepository
}

func NewURLService(r *repository.URLRepository) *urlService {
	return &urlService{r}
}

func (us *urlService) GetURL(ctx context.Context, id, uID string) (*models.URL, error) {
	url, err := us.repo.GetOne(ctx, id, uID)

	if err != nil {
		return nil, models.ErrRecordNotFound
	}

	return url, nil
}
