-- name: GetRefreshToken :one
SELECT * FROM refreshtokens
where $1 = token;


-- name: UpdateRefreshToken :one
UPDATE refreshtokens
SET revoked_at = $1,
        updated_at = $2
where $3 = token
RETURNING *;

-- name: CreateRefreshToken :one
INSERT INTO refreshtokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
       $1, $2, $3, $4, $5, $6
    )
RETURNING *;
