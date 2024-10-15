-- +goose Up
CREATE TABLE users(
    id uuid NOT NULL UNIQUE,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL DEFAULT 'unset'
);

-- +goose Down
DROP TABLE users;
