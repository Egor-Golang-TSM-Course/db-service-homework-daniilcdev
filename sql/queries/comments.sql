-- name: CreateComment :one
INSERT INTO post_comments (id, created_at, post_text, user_id, post_id)
VALUES (
        DEFAULT,
        $1,
        $2,
        $3,
        $4
    )
RETURNING *;

-- name: GetPostComments :many
SELECT posts.* FROM posts
JOIN post_comments ON posts.id == post_comments.post_id
WHERE posts.id == $1
ORDER BY posts.created_at DESC
LIMIT $2;
