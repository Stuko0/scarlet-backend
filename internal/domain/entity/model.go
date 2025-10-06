package entity

type Entity struct {
	EntityId   int64  `json:"entity_id"`
	Name       string `json:"name"`
	EntityType string `json:"entity_type"`
	Address    string `json:"address"`
	Phone      int64 `json:"phone"`
	Email      string `json:"email"`
	ImageUrl   string `json:"image_url"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	Active     bool   `json:"active"`
}