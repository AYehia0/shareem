-- name: CreateShare :one
INSERT INTO share (
  id, url, title, note, ip, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: ListShares :many
SELECT * FROM share ORDER BY created_at DESC;


-- name: CountShares :one
SELECT COUNT(*) FROM share;
