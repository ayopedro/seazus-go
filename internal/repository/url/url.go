package url

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ayopedro/seazus-go/internal/common"
	"github.com/ayopedro/seazus-go/internal/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

type Repository interface {
	GetOne(ctx context.Context, id, uID string) (*models.URL, error)
	GetUserURLs(ctx context.Context, uID string) ([]models.URL, error)
	GetOriginalURL(ctx context.Context, short_url string) (string, error)
	CreateShortURL(ctx context.Context, payload *models.CreateURL, uID string) (string, error)
	UpdateURL(ctx context.Context, id string, payload *models.UpdateURL, uID string) (*models.URL, error)
	DeleteURL(ctx context.Context, id, uID string) error
}

type urlRepository struct {
	client *sql.DB
	logger *zap.Logger
}

func NewRepository(c *sql.DB, l *zap.Logger) Repository {
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
		return nil, err
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
		return nil, err
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

func (ur *urlRepository) CreateShortURL(ctx context.Context, payload *models.CreateURL, uID string) (string, error) {
	var existing string

	err := ur.client.QueryRowContext(
		ctx,
		`SELECT short_url FROM urls WHERE user_id = $1 AND url_address = $2`,
		uID,
		payload.Url,
	).Scan(&existing)

	if err == nil {
		return "", err
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
			return "", err
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

		if pgErr, ok := errors.AsType[*pq.Error](err); ok {
			if pgErr.Code == "23505" {
				continue
			}
		}

		if err == nil {
			return result, nil
		}
		return "", err
	}
	return "", err
}

func (ur *urlRepository) UpdateURL(ctx context.Context, id string, payload *models.UpdateURL, uID string) (*models.URL, error) {
	query := `
		UPDATE urls
		SET
			title = $1,
			description = $2,
			url_address = $3
		WHERE id = $4
		AND user_id = $5
		RETURNING id, title, description, url_address, user_id;
	`

	var updatedURL models.URL
	err := ur.client.QueryRowContext(
		ctx,
		query,
		payload.Identifier,
		payload.Description,
		payload.Url,
		id,
		uID,
	).Scan(
		&updatedURL.Id,
		&updatedURL.Identifier,
		&updatedURL.Description,
		&updatedURL.Url,
		&updatedURL.UserID,
	)

	if err != nil {
		return nil, err
	}

	return &updatedURL, nil
}

func (ur *urlRepository) DeleteURL(ctx context.Context, id, uID string) error {
	query := `
		DELETE FROM urls
		WHERE id = $1
		AND user_id = $2
	`

	res, err := ur.client.ExecContext(ctx, query, id, uID)

	if ra, _ := res.RowsAffected(); ra == 0 {
		return errors.New("Failed to delete resource")
	}

	if err != nil {
		return err
	}
	return nil
}
