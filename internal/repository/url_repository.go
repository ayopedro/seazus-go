package repository

import (
	"context"
	"database/sql"
	"errors"

	appErrors "github.com/ayopedro/seazus-go/internal/common/app_errors"
	"github.com/ayopedro/seazus-go/internal/logger"
	"github.com/ayopedro/seazus-go/internal/models"
)

type urlRepository struct {
	client *sql.DB
}

func NewURLRepository(c *sql.DB) URLRepository {
	return &urlRepository{c}
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
			return nil, appErrors.ErrRecordNotFound
		}
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
		logger.Error(err.Error())
		return nil, appErrors.ErrInternalServerError
	}
	defer rows.Close()

	for rows.Next() {
		url := models.URL{}
		err := rows.Scan(
			&url.Id,
			&url.Identifier,
			&url.Url,
			&url.ShortUrl,
			&url.Description,
			&url.UserID,
			&url.CreatedAt,
			&url.UpdatedAt,
		)
		if err != nil {
			return nil, appErrors.ErrInternalServerError
		}
		urls = append(urls, url)
	}

	if err = rows.Err(); err != nil {
		return nil, appErrors.ErrInternalServerError
	}

	return urls, nil
}
