package comments

import (
	"db-service/internal/database"
	"time"

	"github.com/google/uuid"
)

type comments struct {
	Comments []comment `json:"comments"`
}

type comment struct {
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	Author    author    `json:"author"`
}

type author struct {
	Id uuid.UUID `json:"id"`
}

func databasePostCommentToComment(dbComment *database.PostComment) comment {
	return comment{
		Text:      dbComment.CommentText,
		CreatedAt: dbComment.CreatedAt,

		Author: author{
			Id: dbComment.UserID,
		},
	}
}

func databasePostCommentsToComments(dbComments *[]database.PostComment) *comments {
	r := make([]comment, 0, len(*dbComments))

	for _, dbComment := range *dbComments {
		r = append(r,
			databasePostCommentToComment(&dbComment),
		)
	}

	return &comments{
		Comments: r,
	}
}
