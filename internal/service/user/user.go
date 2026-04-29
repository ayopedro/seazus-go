package user

import (
	"context"

	appErrors "github.com/ayopedro/seazus-go/internal/common/app_errors"
	"github.com/ayopedro/seazus-go/internal/models"
	"github.com/ayopedro/seazus-go/internal/repository"
	"go.uber.org/zap"
)

type Service interface {
	GetUserProfile(ctx context.Context, uID string) (*models.User, error)
	GetUserURLs(ctx context.Context, uID string) ([]models.URL, error)
}

type userService struct {
	userRepo repository.UserRepository
	urlRepo  repository.URLRepository
	logger   *zap.Logger
}

func NewService(
	userRepo repository.UserRepository,
	urlRepo repository.URLRepository,
	logger *zap.Logger,
) Service {
	return &userService{userRepo, urlRepo, logger}
}

func (us *userService) GetUserProfile(ctx context.Context, uID string) (*models.User, error) {
	user, err := us.userRepo.Get(ctx, uID)

	if err != nil {
		return nil, appErrors.ErrNotFound
	}

	return user, nil
}

func (us *userService) GetUserURLs(ctx context.Context, uID string) ([]models.URL, error) {
	urls, err := us.urlRepo.GetUserURLs(ctx, uID)

	if err != nil {
		return nil, appErrors.ErrInternal
	}

	return urls, err
}
