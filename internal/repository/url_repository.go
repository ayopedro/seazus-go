package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ayopedro/seazus-go/cmd/api/dto"
	"github.com/ayopedro/seazus-go/internal/common"
	appErrors "github.com/ayopedro/seazus-go/internal/common/app_errors"
	"github.com/ayopedro/seazus-go/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type urlRepository struct {
	client *sql.DB
	logger *zap.Logger
}

func NewURLRepository(c *sql.DB, l *zap.Logger) URLRepository {
	return &urlRepository{c, l}
}

func (ur *urlRepository) GetOriginalURL(ctx context.Context, short_url string) (string, error) {
	query := `
		SELECT url_address FROM urls
		WHERE short_url = $1;
	`
	var su string
	err := ur.client.QueryRowContext(ctx, query, short_url).Scan(&su)

	if err != nil {
		return "", err
	}

	return su, nil
}

func (ur *urlRepository) GetOne(ctx context.Context, id, uID string) (*models.URL, error) {
	query := `
		SELECT 
			id, 
			title, 
			url_address, 
			short_url, 
			description, 
			created_at, 
			user_id,
			updated_at 
		FROM urls
		WHERE id = $1
		AND user_id = $2;
	`
	url := &models.URL{}
	row := ur.client.QueryRowContext(ctx, query, id, uID)

	err := row.Scan(
		&url.Id,
		&url.Identifier,
		&url.Url,
		&url.ShortUrl,
		&url.Description,
		&url.CreatedAt,
		&url.UserID,
		&url.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, appErrors.MapPostgresError(err)

	}
	return url, nil
}

func (ur *urlRepository) GetUserURLs(ctx context.Context, uID string) ([]models.URL, error) {
	query := `
		SELECT 
			id, 
			title, 
			url_address, 
			short_url, 
			description, 
			user_id,
			created_at, 
			updated_at 
		FROM urls
		WHERE  user_id = $1;
	`
	var urls []models.URL

	rows, err := ur.client.QueryContext(ctx, query, uID)
	if err != nil {
		if ur.logger != nil {
			ur.logger.Error("failed to query user URLs", zap.Error(err))
		}
		return nil, appErrors.ErrInternal
	}
	defer rows.Close()

	for rows.Next() {
		url := models.URL{}
		if err := rows.Scan(
			&url.Id,
			&url.Identifier,
			&url.Url,
			&url.ShortUrl,
			&url.Description,
			&url.UserID,
			&url.CreatedAt,
			&url.UpdatedAt,
		); err != nil {
			return nil, err
		}

		urls = append(urls, url)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func (ur *urlRepository) CreateShortURL(ctx context.Context, payload *dto.CreateURLPayload, uID string) (string, error) {
	var existing string

	err := ur.client.QueryRowContext(
		ctx,
		`SELECT short_url FROM urls WHERE user_id = $1 AND url_address = $2`,
		uID,
		payload.Url,
	).Scan(&existing)

	if err == nil {
		return "", appErrors.ErrConflict
	}

	query := `
		INSERT INTO urls (
			id,
			title,
			description,
			url_address,
			short_url,
			user_id
		)
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING short_url
	`

	for range 5 {
		short_url := common.RandStringRunes(7)
		id, err := uuid.NewV7()
		if err != nil {
			return "", appErrors.ErrInternal
		}

		var result string
		err = ur.client.QueryRowContext(
			ctx,
			query,
			id,
			payload.Identifier,
			payload.Description,
			payload.Url,
			short_url,
			uID,
		).Scan(&result)

		if appErrors.IsUniqueViolation(err) {
			continue
		}

		if err == nil {
			return result, nil
		}
		return "", appErrors.ErrInternal
	}
	return "", appErrors.ErrInternal
}
