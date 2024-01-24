package comments

import (
	"context"
	"db-service/handlers/auth"
	"db-service/internal"
	"db-service/internal/database"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type CommentsStorage struct {
	q *database.Queries
}

func NewStorage(q *database.Queries) *CommentsStorage {
	return &CommentsStorage{
		q: q,
	}
}

func (cs *CommentsStorage) CreateComment(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	postIdParam := chi.URLParam(r, "postId")

	if postId, err := strconv.Atoi(postIdParam); err != nil {
		internal.RespondWithError(w, http.StatusBadRequest, internal.ErrInvalidPostId)
	} else {
		userData := ctx.Value(auth.UserData).(*database.User)

		commentData := commentDto{}
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&commentData)

		comment, err := cs.q.CreateComment(ctx,
			database.CreateCommentParams{
				CommentText: commentData.CommentBody,
				PostID:      int32(postId),
				UserID:      userData.ID,
			})

		switch {
		case err != nil:
			internal.RespondWithError(w, http.StatusInternalServerError, err)
		default:
			internal.RespondWithJSON(
				w, http.StatusCreated,
				databasePostCommentToComment(&comment),
			)
		}
	}
}

func (cs *CommentsStorage) GetAllComments(w http.ResponseWriter, r *http.Request) {
	postIdParam := chi.URLParam(r, "postId")
	if postId, err := strconv.Atoi(postIdParam); err != nil {
		internal.RespondWithError(w, http.StatusBadRequest, internal.ErrInvalidPostId)
	} else {
		comments, err := cs.q.GetPostComments(r.Context(),
			database.GetPostCommentsParams{
				PostID: int32(postId),
				Offset: 0,
				Limit:  10,
			})

		switch {
		case err != nil:
			internal.RespondWithError(w, http.StatusInternalServerError, err)
		default:
			internal.RespondWithJSON(
				w, http.StatusOK,
				databasePostCommentsToComments(&comments),
			)
		}
	}
}
