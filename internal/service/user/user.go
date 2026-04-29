package user

import (
	"context"

	"github.com/ayopedro/seazus-go/internal/models"
	"github.com/ayopedro/seazus-go/internal/repository/url"
	"github.com/ayopedro/seazus-go/internal/repository/user"
	"go.uber.org/zap"
)

type Service interface {
	GetUserProfile(ctx context.Context, uID string) (*models.User, error)
	GetUserURLs(ctx context.Context, uID string) ([]models.URL, error)
}

type userService struct {
	userRepo user.Repository
	urlRepo  url.Repository
	logger   *zap.Logger
}

func NewService(
	userRepo user.Repository,
	urlRepo url.Repository,
	logger *zap.Logger,
) Service {
	return &userService{userRepo, urlRepo, logger}
}

func (us *userService) GetUserProfile(ctx context.Context, uID string) (*models.User, error) {
	user, err := us.userRepo.Get(ctx, uID)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *userService) GetUserURLs(ctx context.Context, uID string) ([]models.URL, error) {
	urls, err := us.urlRepo.GetUserURLs(ctx, uID)

	if err != nil {
		return nil, err
	}

	return urls, nil
}
