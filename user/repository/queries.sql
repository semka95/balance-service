-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;
-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;
-- name: UpdateBalance :one
UPDATE users
SET balance = $2
WHERE id = $1
RETURNING *;
-- name: CreateUser :one
INSERT INTO users(name, email, balance)
VALUES ($1, $2, $3)
RETURNING *;