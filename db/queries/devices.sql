-- name: CreateDevice :one
INSERT INTO devices (name) VALUES ($1) RETURNING id, name, created_at;

-- name: GetDevice :one
SELECT id, name, created_at FROM devices WHERE id = $1;
