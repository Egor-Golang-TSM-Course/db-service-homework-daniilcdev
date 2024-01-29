package comments

import (
	"context"
	"db-service/handlers/auth"
	"db-service/internal"
	"db-service/internal/database"
	"db-service/models"
	"encoding/json"
	"net/http"
)

func (cs *CommentsStorage) CreateComment(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	postId, err := internal.PostIdFromURLParams(r)

	if err != nil {
		internal.RespondWithError(w, http.StatusBadRequest, internal.ErrInvalidPostId)
		return
	}

	userData := ctx.Value(auth.UserData).(*database.User)

	commentData := commentDto{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&commentData)
	if err != nil {
		internal.RespondWithError(w, http.StatusBadRequest, internal.ErrInvalidJson)
		return
	}

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

func (cs *CommentsStorage) GetAllComments(w http.ResponseWriter, r *http.Request) {
	postId, err := internal.PostIdFromURLParams(r)

	if err != nil {
		internal.RespondWithError(w, http.StatusBadRequest, internal.ErrInvalidPostId)
		return
	}

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

func databasePostCommentToComment(dbComment *database.PostComment) models.Comment {
	return models.Comment{
		Text:      dbComment.CommentText,
		CreatedAt: dbComment.CreatedAt,

		Author: models.Author{
			Id: dbComment.UserID,
		},
	}
}

func databasePostCommentsToComments(dbComments *[]database.PostComment) *models.Comments {
	r := make([]models.Comment, 0, len(*dbComments))

	for _, dbComment := range *dbComments {
		r = append(r,
			databasePostCommentToComment(&dbComment),
		)
	}

	return &models.Comments{
		Comments: r,
	}
}
