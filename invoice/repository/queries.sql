-- name: GetInvoiceByID :one
SELECT *
FROM invoices
WHERE id = $1
LIMIT 1;
-- name: GetInvoicesByUserID :many
SELECT *
FROM invoices
WHERE user_id = $1
    AND id > $2
LIMIT $3;
-- name: UpdateStatus :one
UPDATE invoices
SET payment_status = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
-- name: CreateInvoice :one
INSERT INTO invoices(user_id, service_id, order_id, amount)
VALUES ($1, $2, $3, $4)
RETURNING *;