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


-- name: SearchUsers :many
SELECT *
FROM users
WHERE 
    (@username IS NULL OR username ILIKE '%' || @username || '%')
    AND (@email IS NULL OR email ILIKE '%' || @email || '%')
    AND (@full_name IS NULL OR full_name ILIKE '%' || @full_name || '%')
    AND (@bio IS NULL OR bio ILIKE '%' || @bio || '%')
    AND (@created_from IS NULL OR created_at >= @created_from)
    AND (@created_to IS NULL OR created_at <= @created_to)
ORDER BY
    CASE 
        WHEN @sort_by = 'username_asc' THEN username
        WHEN @sort_by = 'username_desc' THEN username
        WHEN @sort_by = 'email_asc' THEN email
        WHEN @sort_by = 'email_desc' THEN email
        WHEN @sort_by = 'created_at_asc' THEN created_at::text
        WHEN @sort_by = 'created_at_desc' THEN created_at::text
        ELSE id::text
    END
    || CASE 
        WHEN @sort_by LIKE '%desc' THEN ' DESC'
        ELSE ' ASC'
    END
LIMIT sqlc.narg('limit_param')::int OFFSET sqlc.narg('offset_param')::int;


-- name: UpdateUser :one
UPDATE users
SET
    username = COALESCE(@username, username),
    email = COALESCE(@email, email),
    password_hash = COALESCE(@password_hash, password_hash),
    full_name = COALESCE(@full_name, full_name),
    bio = COALESCE(@bio, bio),
    updated_at = CURRENT_TIMESTAMP
WHERE id = @id
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;