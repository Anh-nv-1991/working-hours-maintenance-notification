-- name: AddReading :one
INSERT INTO readings (device_id, value, at)
VALUES ($1, $2, COALESCE($3, NOW()))
RETURNING id, device_id, value, at;

-- name: LastReading :one
SELECT id, device_id, value, at
FROM readings
WHERE device_id = $1
ORDER BY at DESC, id DESC
LIMIT 1;
