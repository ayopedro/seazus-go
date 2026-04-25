package service

import (
	"context"

	appErrors "github.com/ayopedro/seazus-go/internal/common/app_errors"
	"github.com/ayopedro/seazus-go/internal/models"
	"github.com/ayopedro/seazus-go/internal/repository"
)

type userService struct {
	userRepo repository.UserRepository
	urlRepo  repository.URLRepository
}

func NewUserService(
	userRepo repository.UserRepository,
	urlRepo repository.URLRepository,
) UserService {
	return &userService{userRepo, urlRepo}
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
