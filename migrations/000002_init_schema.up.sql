CREATE TABLE admins (
    id BIGINT NOT NULL UNIQUE,
    date_until timestamptz NOT NULL
);
