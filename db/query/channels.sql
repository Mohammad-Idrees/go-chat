-- name: CreateChannel :one
INSERT INTO channels (
  name
) VALUES (
  sqlc.arg(name)
)
RETURNING *;


-- name: GetChannels :many
SELECT *
FROM channels;

-- name: GetChannelById :one
SELECT *
FROM channels
where id = sqlc.arg(id);