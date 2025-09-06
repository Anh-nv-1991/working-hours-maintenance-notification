-- name: CreateAlert :one
INSERT INTO alerts (device_id, type, message)
VALUES ($1,$2,$3)
RETURNING *;

-- name: ListOpenAlertsByDevice :many
SELECT * FROM alerts
WHERE device_id = $1 AND resolved = FALSE
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ResolveAlert :one
UPDATE alerts SET
  resolved = TRUE,
  resolved_at = NOW(),
  resolved_by = $2
WHERE id = $1
RETURNING *;