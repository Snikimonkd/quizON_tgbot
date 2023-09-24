CREATE TABLE registrations (
    user_id bigint NOT NULL,
    tg_contact text NOT NULL,
    team_id bigint NOT NULL,
    team_name text NOT NULL,
    captain_name text NOT NULL,
    pnohe text NOT NULL,
    group_name text NOT NULL,
    amount text NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);
