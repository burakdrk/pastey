-- name: CreateSession :one
INSERT INTO sessions (
    id,
    user_id,
    refresh_token,
    user_agent,
    ip_address,
    expires_at
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetSession :one
SELECT *
FROM sessions
WHERE id = $1 LIMIT 1;