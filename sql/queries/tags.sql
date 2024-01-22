-- name: CreateTagForPost :exec
-- INSERT INTO post_tags (id, tag, post_id)
-- VALUES (DEFAULT, $1, $2);

-- name: ListTags :many
-- SELECT DISTINCT tag FROM post_tags
-- LIMIT $1;
