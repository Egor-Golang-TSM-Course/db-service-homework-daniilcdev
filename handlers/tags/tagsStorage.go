package tags

import (
	"context"
	"db-service/internal/database"
)

type TagsQueries interface {
	ExecRaw(ctx context.Context, sql string) error

	GetPost(ctx context.Context, id int32) (database.Post, error)
	UpdatePostTags(ctx context.Context, arg database.UpdatePostTagsParams) (database.Post, error)
	AllTags(ctx context.Context, limit int32) ([]database.Tag, error)
}

type TagsStorage struct {
	q TagsQueries
}

func NewStorage(q TagsQueries) *TagsStorage {
	return &TagsStorage{
		q: q,
	}
}
