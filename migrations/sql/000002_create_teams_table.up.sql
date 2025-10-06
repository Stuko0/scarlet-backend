set search_path to scarlet;
create table teams (
    team_id BIGSERIAL PRIMARY KEY,
    entity_id BIGINT NOT NULL REFERENCES entities(entity_id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    team_type VARCHAR(100) NOT NULL,
    capacity BIGINT NOT NULL,
    fires_attended BIGINT DEFAULT 0,
    team_leader_id BIGINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    active BOOLEAN DEFAULT TRUE
);

CREATE TABLE team_members (
    team_member_id BIGSERIAL PRIMARY KEY,
    team_id BIGINT NOT NULL REFERENCES teams(team_id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL,
    role VARCHAR(100) NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
);

CREATE UNIQUE INDEX uq_team_member_active 
ON team_members (team_id, user_id) 
WHERE (active = TRUE);

CREATE INDEX idx_teams_entity_id ON teams (entity_id);
CREATE INDEX idx_team_members_team_id ON team_members (team_id);
CREATE INDEX idx_team_members_user_id ON team_members (user_id);