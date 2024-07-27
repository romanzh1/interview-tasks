-- name: AddItem :exec
INSERT INTO items (user_id, sku, count)
VALUES ($1, $2, $3);

-- name: GetItemsByUserID :many
SELECT user_id, sku, count
FROM items
WHERE user_id = $1;

