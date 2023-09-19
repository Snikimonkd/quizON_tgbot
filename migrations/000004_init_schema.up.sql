CREATE TABLE registrations_draft (
    user_id bigint NOT NULL UNIQUE,
    game_id bigint,
    team_id bigint,
    team_name text,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);

CREATE SEQUENCE team_id_seq START 1;

CREATE INDEX registrations_draft_user_id_idx ON registrations_draft (user_id);
