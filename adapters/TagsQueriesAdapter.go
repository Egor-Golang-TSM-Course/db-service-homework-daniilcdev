package adapters

import (
	"context"
	"database/sql"
	"db-service/internal/database"
)

type TagsQueriesAdapter struct {
	database.Queries
	db *sql.DB
}

func NewTagsQueriesAdapter(q *database.Queries) *TagsQueriesAdapter {
	return &TagsQueriesAdapter{
		Queries: *q,
	}
}
func (tqa *TagsQueriesAdapter) WithDb(db *sql.DB) *TagsQueriesAdapter {
	tqa.db = db
	return tqa
}

func (tqa *TagsQueriesAdapter) ExecRaw(ctx context.Context, sql string) error {
	_, err := tqa.db.ExecContext(ctx, sql)
	return err
}
