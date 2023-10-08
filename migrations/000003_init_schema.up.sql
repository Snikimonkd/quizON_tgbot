CREATE TABLE registrations (
    tg_contact text NOT NULL,
    team_id text,
    team_name text NOT NULL,
    captain_name text NOT NULL,
    phone text NOT NULL,
    group_name text NOT NULL,
    amount text NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);
