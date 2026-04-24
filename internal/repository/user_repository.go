package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/ayopedro/seazus-go/internal/models"
)

type UserRepository struct {
	client *sql.DB
}

func NewUserRepository(c *sql.DB) *UserRepository {
	return &UserRepository{c}
}

func (ur *UserRepository) Create(ctx context.Context, u *models.User) error {
	query := "INSERT INTO users (id, first_name, last_name, email, password) VALUES ($1, $2, $3, $4, $5);"

	_, err := ur.client.ExecContext(ctx, query, u.Id, u.FirstName, u.LastName, u.Email, u.Password)
	if err != nil {
		if strings.Contains(err.Error(), models.Conflict) {
			return models.ErrDuplicateEmail
		}
		return err
	}
	return nil
}

func (ur *UserRepository) Get(ctx context.Context, uId string) (*models.User, error) {
	query := `
		SELECT id, first_name, last_name, email, is_verified, created_at, updated_at
		FROM users
		WHERE email = $1;
		`
	user := &models.User{}
	row := ur.client.QueryRowContext(ctx, query, uId)

	err := row.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.IsVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) GetWithEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, first_name, last_name, email, password, is_verified, created_at, updated_at
		FROM users
		WHERE email = $1;
		`
	user := &models.User{}
	row := ur.client.QueryRowContext(ctx, query, strings.ToLower(email))

	err := row.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.IsVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}
