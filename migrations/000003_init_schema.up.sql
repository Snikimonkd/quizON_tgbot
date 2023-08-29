CREATE TABLE registrations (
    game_id BIGINT NOT NULL REFERENCES games(id),
    team_id text NOT NULL,
    team_name text NOT NULL,
    created_at timestamptz NOT NULL
);

CREATE UNIQUE INDEX game_id_team_id_uniq ON registrations (game_id, team_id);
