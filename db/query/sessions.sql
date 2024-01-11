-- name: CreateSession :one
INSERT INTO sessions (
    id,
    email,
    user_agent,
    client_ip,
    refresh_token,
    expires_at
) VALUES (
    sqlc.arg(id), sqlc.arg(email), sqlc.arg(user_agent), sqlc.arg(client_ip), sqlc.arg(refresh_token), sqlc.arg(expires_at)
) RETURNING *;

-- name: UpdateSession :one
UPDATE sessions
SET
    is_blocked = COALESCE(sqlc.narg(is_blocked), is_blocked),
    is_logged_out = COALESCE(sqlc.narg(is_logged_out),is_logged_out)
WHERE
    id = sqlc.arg(id)
RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions 
WHERE id = sqlc.arg(id);