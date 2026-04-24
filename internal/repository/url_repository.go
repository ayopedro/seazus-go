package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ayopedro/seazus-go/internal/logger"
	"github.com/ayopedro/seazus-go/internal/models"
)

type URLRepository struct {
	client *sql.DB
}

func NewURLRepository(c *sql.DB) *URLRepository {
	return &URLRepository{c}
}

func (ur *URLRepository) GetOne(ctx context.Context, id, uID string) (*models.URL, error) {
	logger.Debug(uID)
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
		logger.Debug(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrRecordNotFound
		}
		return nil, err
	}
	return url, nil
}
