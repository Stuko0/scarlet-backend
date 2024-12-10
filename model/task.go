package model

type Task struct {
	ID          int64  `json:"id"`
	TeamID      int64  `json:"team_id"`
	FireID      string `json:"fire_id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Active      bool   `json:"active"`
	CreatedAt   string `json:"created_at"`
	CreatedBy   int64  `json:"created_by"`
}
