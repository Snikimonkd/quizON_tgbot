CREATE TABLE registrations_draft (
    user_id bigint NOT NULL UNIQUE,
    tg_contact text NOT NULL,
    team_id text,
    team_name text,
    captain_name text,
    group_name text,
    phone text,
    amount text,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);
