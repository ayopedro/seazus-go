package repository

import (
	"context"
	"database/sql"
	"strings"

	appErrors "github.com/ayopedro/seazus-go/internal/common/app_errors"
	"github.com/ayopedro/seazus-go/internal/models"
	"go.uber.org/zap"
)

type userRepository struct {
	client *sql.DB
	logger *zap.Logger
}

func NewUserRepository(c *sql.DB, logger *zap.Logger) UserRepository {
	return &userRepository{client: c, logger: logger}
}

func (ur *userRepository) Create(ctx context.Context, u *models.User) error {
	query := "INSERT INTO users (id, first_name, last_name, email, password) VALUES ($1, $2, $3, $4, $5);"

	_, err := ur.client.ExecContext(ctx, query, u.Id, u.FirstName, u.LastName, u.Email, u.Password)
	if err != nil {
		return appErrors.MapPostgresError(err)
	}
	return nil
}

func (ur *userRepository) Get(ctx context.Context, uId string) (*models.User, error) {
	query := `
		SELECT id, first_name, last_name, email, is_verified, created_at, updated_at
		FROM users
		WHERE id = $1;
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
		return nil, appErrors.MapPostgresError(err)
	}

	return user, nil
}

func (ur *userRepository) GetWithEmail(ctx context.Context, email string) (*models.User, error) {
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
		return nil, appErrors.MapPostgresError(err)
	}

	return user, nil
}
