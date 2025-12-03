-- name: GetUserByEmailOrUsername :one
SELECT id, username, email, password, role_id FROM users WHERE email = $1 OR username = $1;

-- name: CreateUser :one
INSERT INTO users (id, username, email, password, role_id) VALUES ($1, $2, $3, $4, $5)
RETURNING id, username, email, role_id, created_at, updated_at;