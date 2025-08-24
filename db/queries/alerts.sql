-- name: CreateAlert :one
INSERT INTO alerts (device_id, reading_id, level, message, is_serviced, created_at)
VALUES ($1, $2, $3, $4, FALSE, COALESCE($5, NOW()))
RETURNING id, device_id, reading_id, level, message, is_serviced, created_at, serviced_at;

-- name: ResolveAlert :one
UPDATE alerts
SET is_serviced = TRUE, serviced_at = COALESCE($2, NOW())
WHERE id = $1
RETURNING id, device_id, reading_id, level, message, is_serviced, created_at, serviced_at;

-- name: ListOpenAlertsByDevice :many
SELECT id, device_id, reading_id, level, message, is_serviced, created_at, serviced_at
FROM alerts
WHERE device_id = $1 AND is_serviced = FALSE
ORDER BY created_at DESC;

-- name: CheckThresholdBreach :one
SELECT
  CASE
    WHEN lr.value < p.threshold_min THEN 'LOW'
    WHEN lr.value > p.threshold_max THEN 'HIGH'
    ELSE 'OK'
  END AS status,
  lr.id    AS reading_id,
  lr.value AS reading_value
FROM plans p
JOIN LATERAL (
  SELECT id, device_id, value, at
  FROM readings r
  WHERE r.device_id = p.device_id
  ORDER BY at DESC, id DESC
  LIMIT 1
) lr ON TRUE
WHERE p.device_id = $1;
