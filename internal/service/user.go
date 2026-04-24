package service

import (
	"context"
	"strings"
	"time"

	"github.com/ayopedro/seazus-go/internal/logger"
	"github.com/ayopedro/seazus-go/internal/models"
	"github.com/ayopedro/seazus-go/internal/repository"
	"github.com/ayopedro/seazus-go/internal/utils"
	"github.com/google/uuid"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo}
}

func (us *UserService) CreateUser(ctx context.Context, u *models.User) error {
	uID, _ := uuid.NewV7()

	hash, err := utils.HashPassword(u.Password)
	if err != nil {
		logger.Error("Error hashing password")
		return models.ErrInternalServerError
	}

	user := &models.User{
		Id:        uID.String(),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     strings.ToLower(u.Email),
		Password:  hash,
		CreatedAt: time.Now(),
	}

	err = us.repo.Create(ctx, user)
	return err
}

func (us *UserService) GetUserProfile(ctx context.Context, uID string) (*models.User, error) {
	user, err := us.repo.Get(ctx, uID)

	if err != nil {
		return nil, err
	}

	return user, nil
}
