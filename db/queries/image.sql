-- name: CreateImage :one
INSERT INTO images (id, original_name, size, mime_type, path)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, original_name, size, mime_type, path, created_at;

-- name: GetImage :one
SELECT id, original_name, size, mime_type, path, created_at
FROM images
WHERE id = $1;

-- name: DeleteImage :exec
DELETE FROM images
WHERE id = $1;