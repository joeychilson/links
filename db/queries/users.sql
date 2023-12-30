-- name: CreateUser :one
INSERT INTO users (avatar, email, username, password) VALUES ($1, $2, $3, $4) RETURNING id;

-- name: UserByID :one
SELECT id, avatar, username, email, confirmed_at, created_at, updated_at FROM users WHERE id = $1;

-- name: EmailExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);

-- name: UsernameExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE username = $1);

-- name: CreateUserToken :one
INSERT INTO user_tokens (user_id, token, context) VALUES ($1, $2, $3) RETURNING token;

-- name: DeleteUserToken :exec
DELETE FROM user_tokens WHERE token = $1 AND context = $2;

-- name: UserIDByToken :one
SELECT user_id FROM user_tokens WHERE token = $1 AND context = $2;