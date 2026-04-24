package models

type CreateURLPayload struct {
	Identifier  string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url_address"`
}

type URL struct {
	Id          string   `json:"id"`
	Identifier  string   `json:"title"`
	Description string   `json:"description"`
	Url         string   `json:"url_address"`
	ShortUrl    string   `json:"short_url"`
	CreatedAt   JSONTime `json:"created_at"`
	UpdatedAt   JSONTime `json:"updated_at"`
	UserID      string   `json:"user_id"`
}
