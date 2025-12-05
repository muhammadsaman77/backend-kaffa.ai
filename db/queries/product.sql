-- name: CreateProduct :one
INSERT INTO products (id, store_id, image_id, name, description, price, is_available)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, store_id, image_id, name, description, price, is_available, created_at, updated_at;

-- name: GetListProductsByStoreId :many
SELECT 
    p.id AS id,
    p.name AS name,
    p.description AS description,
    p.price AS price,
    p.is_available AS is_available,
    i.path AS path,
    p.created_at AS created_at,
    p.updated_at AS updated_at
FROM products p
LEFT JOIN images i ON p.image_id = i.id
WHERE p.store_id = $1;


-- name: GetProductById :one
SELECT 
    p.id AS id,
    p.store_id AS store_id,
    p.image_id AS image_id,
    p.name AS name,
    p.description AS description,
    p.price AS price,
    p.is_available AS is_available,
    i.path AS path,
    p.created_at AS created_at,
    p.updated_at AS updated_at
FROM products p
LEFT JOIN images i ON p.image_id = i.id
WHERE p.id = $1;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;
