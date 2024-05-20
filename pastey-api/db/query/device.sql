-- name: CreateDevice :one
INSERT INTO devices (user_id, device_name, public_key)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetDeviceById :one
SELECT *
FROM devices
WHERE id = $1
LIMIT 1;

-- name: ListUserDevices :many
SELECT *
FROM devices
WHERE user_id = $1;

-- name: UpdateDevice :one
UPDATE devices
SET
    device_name = COALESCE(sqlc.narg(device_name), device_name),
    public_key = COALESCE(sqlc.narg(public_key), public_key)
WHERE
    id = sqlc.arg(id)
RETURNING *;

-- name: DeleteDevice :exec
DELETE FROM devices
WHERE id = $1;
