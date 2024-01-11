-- name: CreateMembership :one
INSERT INTO memberships (
  user_id, channel_id
) VALUES (
  sqlc.arg(user_id), sqlc.arg(channel_id) 
)
RETURNING *;


-- name: GetMemberships :many
SELECT *
FROM memberships;

-- name: GetMembershipsByUserId :many
SELECT *
FROM memberships
where user_id = sqlc.arg(user_id);

-- name: GetMembershipsByChannelId :many
SELECT *
FROM memberships
where channel_id = sqlc.arg(channel_id);