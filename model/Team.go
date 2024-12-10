package model

type Member struct {
	UserId    int64  `json:"user_id"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type Team struct {
	Id            int64    `json:"id"`
	Name          string   `json:"name"`
	Members       []Member `json:"members"`
	FiresAttended int64    `json:"fires_attended"`
	Active        bool     `json:"active"`
	CreatedAt     string   `json:"created_at"`
	CreatedBy     int64    `json:"created_by"`
}
