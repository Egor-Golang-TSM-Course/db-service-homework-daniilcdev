package comments

import "db-service/internal/database"

type CommentsStorage struct {
	q *database.Queries
}

func NewStorage(q *database.Queries) *CommentsStorage {
	return &CommentsStorage{
		q: q,
	}
}
