// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: tags.sql

package database

import (
	"context"
)

const createTagForPost = `-- name: CreateTagForPost :exec
INSERT INTO post_tags (id, tag, post_id)
VALUES (DEFAULT, $1, $2)
`

type CreateTagForPostParams struct {
	Tag    string
	PostID int32
}

func (q *Queries) CreateTagForPost(ctx context.Context, arg CreateTagForPostParams) error {
	_, err := q.db.ExecContext(ctx, createTagForPost, arg.Tag, arg.PostID)
	return err
}

const listTags = `-- name: ListTags :many
SELECT DISTINCT tag FROM post_tags
LIMIT $1
`

func (q *Queries) ListTags(ctx context.Context, limit int32) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, listTags, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		items = append(items, tag)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
