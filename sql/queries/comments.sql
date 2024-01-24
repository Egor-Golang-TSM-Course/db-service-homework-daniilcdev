-- name: CreateComment :one
INSERT INTO post_comments (id, created_at, comment_text, user_id, post_id)
VALUES (
        DEFAULT,
        NOW(),
        $1,
        $2,
        $3
    )
RETURNING *;

-- name: GetPostComments :many
SELECT * FROM post_comments
WHERE post_id = $1
ORDER BY created_at DESC
OFFSET $2
LIMIT $3;
