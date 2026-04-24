-- +goose Up
CREATE TABLE IF NOT EXISTS urls (
    id UUID PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    url_address TEXT NOT NULL,
    short_url VARCHAR(10),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS short_url_unique ON urls (short_url);

-- +goose Down
DROP TABLE IF EXISTS urls;
