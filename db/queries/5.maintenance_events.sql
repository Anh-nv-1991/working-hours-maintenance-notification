INSERT INTO maintenance_events (device_id, at, interval, notes, performed_by, cost)
VALUES ($1,$2,$3,$4,$5,$6)
RETURNING *;

-- name: ListMaintenanceByDevice :many
SELECT * FROM maintenance_events
WHERE device_id = $1
ORDER BY at DESC
LIMIT $2 OFFSET $3;

-- name: DeleteMaintenanceEvent :exec
DELETE FROM maintenance_events WHERE id = $1;