-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, content, user_id)
VALUES (
        DEFAULT,
        $1,
        $2,
        $3,
        $4,
        $5
    )
RETURNING *;

-- name: GetPostsByUser :many
SELECT posts.* FROM posts
JOIN post_tags ON post_tags.post_id == posts.id
WHERE posts.user_id = $1
ORDER BY posts.created_at DESC
LIMIT $2;
