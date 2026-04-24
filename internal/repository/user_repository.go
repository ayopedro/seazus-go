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
	query := "SELECT * FROM users WHERE id = $1;"
	user := &models.User{}
	row := ur.client.QueryRowContext(ctx, query, uId)

	err := row.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
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
