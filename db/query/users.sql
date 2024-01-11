-- name: CreateUser :one
INSERT INTO users (
  username, email, hashed_password, phone
) VALUES (
  sqlc.arg(username), sqlc.arg(email), sqlc.arg(hashed_password), sqlc.arg(phone)
)
RETURNING *;


-- name: GetUserByEmail :one
SELECT *
FROM users
where email = sqlc.arg(email);

-- name: GetUserById :one
SELECT *
FROM users
where id = sqlc.arg(id);

-- name: GetUsers :many
SELECT *
FROM users;