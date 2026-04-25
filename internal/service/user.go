package service

import (
	"context"

	appErrors "github.com/ayopedro/seazus-go/internal/common/app_errors"
	"github.com/ayopedro/seazus-go/internal/models"
	"github.com/ayopedro/seazus-go/internal/repository"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

func (us *userService) GetUserProfile(ctx context.Context, uID string) (*models.User, error) {
	user, err := us.repo.Get(ctx, uID)

	if err != nil {
		return nil, appErrors.ErrUserNotFound
	}

	return user, nil
}
