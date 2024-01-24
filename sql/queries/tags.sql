-- name: UpdatePostTags :one
UPDATE posts
SET tags = $3,
    updated_at = NOW()
WHERE posts.id = $1 AND posts.user_id = $2
RETURNING *;

-- name: AddTag :exec
INSERT INTO tags (id, tag, created_at) 
VALUES (DEFAULT, UNNEST($1), NOW()) ON CONFLICT DO NOTHING;
 

-- name: AllTags :many
SELECT * FROM tags
LIMIT $1;