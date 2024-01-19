-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, email, access_token)
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6
    )
RETURNING *;

-- name: AuthorizeUser :one
SELECT * FROM users
WHERE access_token = $1
LIMIT 1;
