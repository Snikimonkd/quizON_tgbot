CREATE TABLE games (
    id BIGSERIAL PRIMARY KEY,
    created_at timestamptz NOT NULL,
    udpated_at timestamptz NOT NULL,
    location text NOT NULL,
    date timestamptz NOT NULL,
    description text
);
