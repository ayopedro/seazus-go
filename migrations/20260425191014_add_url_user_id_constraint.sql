-- +goose Up
ALTER TABLE urls
ADD CONSTRAINT urls_user_url_unique UNIQUE (user_id, url_address);
