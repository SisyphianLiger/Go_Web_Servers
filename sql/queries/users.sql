-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
        gen_random_uuid(), NOW(), NOW(), $1, $2
    )
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1;

-- name: ResetAllUsers :exec
DELETE FROM users;


-- name: UpdateUserInfo :one
UPDATE users
SET email = $1,
        hashed_password = $2
where $3 = id
RETURNING *;
