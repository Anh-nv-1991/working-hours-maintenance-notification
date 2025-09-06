-- name: CreateDevice :one
INSERT INTO devices (
  serial_number, name, model, manufacturer, year_of_manufacture,
  commission_date, total_working_hour, after_overhaul_working_hour,
  status, last_service_at, location, plan_id, created_at, updated_at
) VALUES (
  $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,NOW(),NOW()
) RETURNING *;

-- name: GetDevice :one
SELECT * FROM devices WHERE id = $1 LIMIT 1;

-- name: ListDevices :many
SELECT * FROM devices
WHERE deleted_at IS NULL
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateDeviceBasic :one
UPDATE devices SET
  name = $2,
  status = $3,
  location = $4,
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateDevicePlan :one
UPDATE devices SET
  plan_id = $2,
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: SoftDeleteDevice :exec
UPDATE devices SET deleted_at = NOW() WHERE id = $1;
