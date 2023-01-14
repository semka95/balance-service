-- name: GetTransferByID :one
SELECT *
FROM transfers
WHERE id = $1
LIMIT 1;
-- name: GetInboundTransfers :many
SELECT *
FROM transfers
WHERE to_user_id = $1
    AND id > $2
LIMIT $3;
-- name: GetOutboundTransfers :many
SELECT *
FROM transfers
WHERE from_user_id = $1
    AND id > $2
LIMIT $3;
-- name: GetTransfersBetweenUsers :many
SELECT *
FROM transfers
WHERE from_user_id = $1
    AND to_user_id = $2
    AND id > $3
LIMIT $4;
-- name: CreateTransfer :one
INSERT INTO transfers(from_user_id, to_user_id, amount)
VALUES ($1, $2, $3)
RETURNING *;