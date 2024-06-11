-- name: CreateEntry :one
INSERT INTO clipboard_entries (entry_id, user_id, from_device_id, to_device_id, encrypted_data, created_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetEntriesForDevice :many
SELECT c.*, d.device_name as from_device_name
FROM clipboard_entries c
JOIN devices d ON c.from_device_id = d.id
WHERE c.to_device_id = $1
ORDER BY c.created_at DESC;

-- name: GetEntryByEntryId :many
SELECT *
FROM clipboard_entries
WHERE entry_id = $1
ORDER BY created_at DESC;

-- name: GetEntryByUser :many
SELECT *
FROM clipboard_entries
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetEntryByUserForUpdate :many
SELECT *
FROM clipboard_entries
WHERE user_id = $1
ORDER BY created_at DESC
FOR NO KEY UPDATE;

-- name: DeleteEntry :exec
DELETE FROM clipboard_entries
WHERE entry_id = $1;
