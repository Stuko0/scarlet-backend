package teamdata

import (
	"database/sql"
	"scarlet_backend/model"
)

type MemberDataPostgres struct{ db *sql.DB }

func NewMemberDataPostgres(db *sql.DB) *MemberDataPostgres { return &MemberDataPostgres{db: db} }

func (m *MemberDataPostgres) GetMembers() ([]model.Member, error) {
	rows, err := m.db.Query("SELECT memberId, team_id, user_id, roleMember, created_at FROM team_data.members")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var members []model.Member
	for rows.Next() {
		var member model.Member
		if err = rows.Scan(&member.UserId, &member.Role); err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, nil
}

func (m *MemberDataPostgres) AddMember(member model.Member) error {
	query := `INSERT INTO team_data.members (team_id, user_id, roleMember, created_at) VALUES ($1, $2, $3, $4)`
	_, err := m.db.Exec(query, member.UserId, member.Role, member.CreatedAt)
	return err
}
