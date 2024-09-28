-- name: CreateUser :one
INSERT INTO users (
    username, email, password_hash, full_name, bio
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET 
    username = COALESCE($2, username),
    email = COALESCE($3, email),
    password_hash = COALESCE($4, password_hash),
    full_name = COALESCE($5, full_name),
    bio = COALESCE($6, bio),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;