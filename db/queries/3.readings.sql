-- name: CreateReading :one
INSERT INTO readings (device_id, at, hours_delta, location, operator_id)
VALUES ($1,$2,$3,$4,$5)
RETURNING *;

-- name: ListReadingsByDevice :many
SELECT * FROM readings
WHERE device_id = $1
ORDER BY at DESC
LIMIT $2 OFFSET $3;

-- name: GetLastReading :one
SELECT * FROM readings
WHERE device_id = $1
ORDER BY at DESC
LIMIT 1;

-- name: DeleteReading :exec
DELETE FROM readings WHERE id = $1;
