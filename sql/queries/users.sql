-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, email, pwd_hash)
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6
    )
RETURNING *;

-- name: UserByEmail :one
SELECT * FROM users
WHERE email = $1
LIMIT 1;

-- name: UserById :one
SELECT * FROM users
WHERE id = $1
LIMIT 1;