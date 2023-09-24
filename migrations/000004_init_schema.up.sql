CREATE TABLE registrations_draft (
    user_id bigint NOT NULL,
    tg_contact text NOT NULL,
    team_id bigint NOT NULL,
    team_name text,
    captain_name text,
    group_name text,
    phone text,
    amount text,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);

CREATE SEQUENCE team_id_seq START 1;
