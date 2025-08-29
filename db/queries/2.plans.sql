-- name: CreatePlan :one
INSERT INTO plans (name, interval_hours, description, created_at, updated_at)
VALUES ($1,$2,$3,NOW(),NOW()) RETURNING *;

-- name: GetPlan :one
SELECT * FROM plans WHERE id = $1 LIMIT 1;

-- name: ListPlans :many
SELECT * FROM plans ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdatePlan :one
UPDATE plans SET
  name = $2,
  interval_hours = $3,
  description = $4,
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeletePlan :exec
DELETE FROM plans WHERE id = $1;
