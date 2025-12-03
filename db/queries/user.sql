-- name: GetUserByEmailOrUsername :one
SELECT id, username, email, password, role_id FROM users WHERE email = $1 OR username = $1;