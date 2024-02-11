package comments

import (
	"context"
	"db-service/handlers"
	"db-service/handlers/auth"
	"db-service/internal"
	"db-service/internal/database"
	"db-service/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func (cs *CommentsStorage) CreateComment(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	postId, err := internal.PostIdFromURLParams(r)

	if err != nil {
		return &handlers.HttpError{Code: http.StatusBadRequest, Err: internal.ErrInvalidPostId}
	}

	userData := ctx.Value(auth.UserData).(*database.User)

	var commentData commentDto
	if err = json.NewDecoder(r.Body).Decode(&commentData); err != nil {
		return &handlers.HttpError{Code: http.StatusBadRequest, Err: internal.ErrInvalidJson}
	}

	comment, err := cs.q.CreateComment(ctx,
		database.CreateCommentParams{
			CommentText: commentData.CommentBody,
			PostID:      int32(postId),
			UserID:      userData.ID,
		})

	if err != nil {
		return &handlers.HttpError{Code: http.StatusInternalServerError, Err: err}
	}

	fmt.Fprint(w, databasePostCommentToComment(&comment))
	internal.SetHeaders(w, http.StatusCreated)
	return nil
}

func (cs *CommentsStorage) GetAllComments(w http.ResponseWriter, r *http.Request) error {
	postId, err := internal.PostIdFromURLParams(r)

	if err != nil {
		return &handlers.HttpError{Code: http.StatusBadRequest, Err: internal.ErrInvalidPostId}
	}

	comments, err := cs.q.GetPostComments(r.Context(),
		database.GetPostCommentsParams{
			PostID: int32(postId),
			Offset: 0,
			Limit:  10,
		})

	if err != nil {
		return &handlers.HttpError{Code: http.StatusInternalServerError, Err: err}
	}

	fmt.Fprint(w, databasePostCommentsToComments(&comments))
	internal.SetHeaders(w, http.StatusOK)

	return nil
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
