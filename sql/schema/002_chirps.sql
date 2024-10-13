-- +goose Up
CREATE TABLE chirps(
    id uuid NOT NULL UNIQUE,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    body TEXT NOT NULL,
    user_id uuid NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE chirps;
