-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, content, user_id)
VALUES (
        DEFAULT,
        NOW(),
        NOW(),
        $1,
        $2,
        $3
    )
RETURNING *;

-- name: GetPost :one
SELECT *
FROM posts
WHERE posts.id = $1;

-- name: GetPosts :many
SELECT * FROM posts
ORDER BY posts.created_at DESC
OFFSET $1
LIMIT $2;

-- name: GetPostsByUser :many
SELECT posts.* FROM posts
WHERE posts.user_id = $1
ORDER BY posts.created_at DESC
LIMIT $2;

-- name: UpdatePost :one
UPDATE posts
SET updated_at = NOW(),
 title = $3,
 content = $4
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: UpdatePostTitle :one
UPDATE posts
SET updated_at = NOW(),
 title = $3
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: UpdatePostContent :one
UPDATE posts
SET updated_at = NOW(),
 content = $3
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1 AND user_id = $2;

