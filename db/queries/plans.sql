-- name: UpsertPlan :one
INSERT INTO plans (device_id, threshold_min, threshold_max, created_at)
VALUES ($1, $2, $3, COALESCE($4, NOW()))
ON CONFLICT (device_id) DO UPDATE
SET threshold_min = EXCLUDED.threshold_min,
    threshold_max = EXCLUDED.threshold_max
RETURNING device_id, threshold_min, threshold_max, created_at;

-- name: GetPlanForDevice :one
SELECT device_id, threshold_min, threshold_max, created_at
FROM plans
WHERE device_id = $1;
