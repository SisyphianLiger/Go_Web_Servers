-- +goose Up
CREATE TABLE refreshtokens(
    token TEXT NOT NULL PRIMARY KEY,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    user_id uuid NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    expires_at timestamptz NOT NULL,
    revoked_at timestamptz
);

-- +goose Down
DROP TABLE refreshtokens; 

