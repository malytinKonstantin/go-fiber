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
    (@username::text IS NULL OR username ILIKE '%' || @username::text || '%')
    AND (@email::text IS NULL OR email ILIKE '%' || @email::text || '%')
    AND (@full_name::text IS NULL OR full_name ILIKE '%' || @full_name::text || '%')
    AND (@bio::text IS NULL OR bio ILIKE '%' || @bio::text || '%')
    AND (@created_from::timestamptz IS NULL OR created_at >= @created_from::timestamptz)
    AND (@created_to::timestamptz IS NULL OR created_at <= @created_to::timestamptz)
ORDER BY
    CASE 
        WHEN @sort_by::text = 'username_asc' THEN username
        WHEN @sort_by::text = 'email_asc' THEN email
        WHEN @sort_by::text = 'created_at_asc' THEN NULL -- We’ll handle this separately
        WHEN @sort_by::text = 'id_asc' THEN NULL -- We’ll handle this separately
    END ASC,
    CASE 
        WHEN @sort_by::text = 'username_desc' THEN username
        WHEN @sort_by::text = 'email_desc' THEN email
        WHEN @sort_by::text = 'created_at_desc' THEN NULL -- We’ll handle this separately
        WHEN @sort_by::text = 'id_desc' THEN NULL -- We’ll handle this separately
    END DESC,
    CASE WHEN @sort_by::text = 'created_at_asc' THEN created_at END ASC,
    CASE WHEN @sort_by::text = 'id_asc' THEN id END ASC,
    CASE WHEN @sort_by::text = 'created_at_desc' THEN created_at END DESC,
    CASE WHEN @sort_by::text = 'id_desc' THEN id END DESC,
    id ASC -- Always fallback sort by id
LIMIT sqlc.narg('limit_param')::int
OFFSET sqlc.narg('offset_param')::int;


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