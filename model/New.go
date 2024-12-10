package model

type New struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	URL         string `json:"url"`
	Active      bool   `json:"active"`
	CreatedAt   string `json:"created_at"`
	CreatedBy   int64  `json:"created_by"`
}
