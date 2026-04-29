package dto

type CreateURLPayload struct {
	Identifier  string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url_address"`
}
