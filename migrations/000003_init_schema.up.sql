CREATE TABLE registrations (
    game_id bigint NOT NULL,
    team_id bigint NOT NULL,
    team_name text NOT NULL,
    user_id bigint NOT NULL,
    created_at timestamptz NOT NULL,
    udpated_at timestamptz NOT NULL
);

CREATE UNIQUE INDEX game_id_team_id_uniq ON registrations (game_id, team_id);

