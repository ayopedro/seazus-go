package repository

import (
	"context"
	"database/sql"
	"errors"

	appErrors "github.com/ayopedro/seazus-go/internal/common/app_errors"
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
