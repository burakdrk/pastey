-- name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;

-- name: GetUserById :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
    password_hash = COALESCE(sqlc.narg(password_hash), password_hash),
    ispremium = COALESCE(sqlc.narg(ispremium), ispremium),
    isemailverified = COALESCE(sqlc.narg(isemailverified), isemailverified),
    email = COALESCE(sqlc.narg(email), email)
WHERE
    id = sqlc.arg(id)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
