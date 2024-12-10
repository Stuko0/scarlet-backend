package teamdata

import (
	"database/sql"
	"scarlet_backend/model"
)

type TeamDataPostgres struct{ db *sql.DB }

func NewTeamDataPostgres(db *sql.DB) *TeamDataPostgres { return &TeamDataPostgres{db: db} }

func (t *TeamDataPostgres) GetTeams() ([]model.Team, error) {
	rows, err := t.db.Query("SELECT teamId, teamName, fires_attended, active, created_at, created_by FROM team_data.teams")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var teams []model.Team
	for rows.Next() {
		var team model.Team
		if err = rows.Scan(&team.Id, &team.Name, &team.FiresAttended, &team.Active, &team.CreatedAt, &team.CreatedBy); err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	return teams, nil
}

func (t *TeamDataPostgres) AddTeam(team model.Team) error {
	query := `INSERT INTO team_data.teams (teamName, fires_attended, active, created_at, created_by) VALUES ($1, $2, $3, $4, $5)`
	_, err := t.db.Exec(query, team.Name, team.FiresAttended, team.Active, team.CreatedAt, team.CreatedBy)
	return err
}
