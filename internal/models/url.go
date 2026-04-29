package models

import "time"

type URL struct {
	Id          string
	Identifier  string
	Description string
	Url         string
	ShortUrl    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserID      string
}

type CreateURL struct {
	Identifier  string
	Description string
	Url         string
}
