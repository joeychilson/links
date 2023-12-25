-- name: CreateUser :exec
INSERT INTO users (username, email, password) VALUES ($1, $2, $3);

-- name: UpdateUser :exec
UPDATE users SET email = $1, username = $2 WHERE id = $3;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: FindUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: CheckEmailExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);

-- name: CheckUsernameExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE username = $1);

-- name: CreateUserToken :exec
INSERT INTO user_tokens (user_id, token, context, expires_at) VALUES ($1, $2, $3, $4);

-- name: UpdateUserToken :exec
UPDATE user_tokens SET expires_at = $1 WHERE id = $2;

-- name: DeleteUserToken :exec
DELETE FROM user_tokens WHERE id = $1;

-- name: FindTokensByUserID :many
SELECT * FROM user_tokens WHERE user_id = $1;

-- name: DeleteExpiredTokens :exec
DELETE FROM user_tokens WHERE expires_at < CURRENT_TIMESTAMP;