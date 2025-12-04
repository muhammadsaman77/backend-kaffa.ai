-- name: CreateProduct :one
INSERT INTO products (id, store_id, name, description, price, is_available)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, store_id, name, description, price, is_available, created_at, updated_at;