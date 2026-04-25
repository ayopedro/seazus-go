package service

import (
	"context"

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
		return nil, models.ErrUserNotFound
	}

	return user, nil
}
