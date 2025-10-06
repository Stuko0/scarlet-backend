package team

type Team struct {
	TeamId    int64  `json:"team_id"`
	EntityId    int64  `json:"entity_id"`
	Name      string `db:"name"`
	TeamType string `db:"team_type"`
	Capacity int64  `db:"capacity"`
	FiresAttended int64  `db:"fires_attended"`
	TeamLeaderId int64  `db:"team_leader_id"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
	Active    bool   `db:"active"`
}

type TeamMember struct {
	TeamMemberId int64  `json:"team_member_id"`
	TeamId    int64  `json:"team_id"`
	UserId    int64  `json:"user_id"`
	Role      string `db:"role"`
	Active	bool   `db:"active"`
	JoinedAt  string `db:"joined_at"`
}