package posts

import (
	"context"
	"database/sql"
	"db-service/handlers/auth"
	"db-service/internal"
	"db-service/internal/database"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

var errEmptyTitle = errors.New("restriction violation: title can't be empty")
var errMissingPostContent = errors.New("missing post data")
var errNotFound = errors.New("post not found")

type PostsStorage struct {
	q *database.Queries
}

func NewStorage(q *database.Queries) *PostsStorage {
	return &PostsStorage{
		q: q,
	}
}

func (s *PostsStorage) CreatePost(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	params := createPostRequestParams{}
	err := decoder.Decode(&params)
	var post database.Post

	switch {
	case err == io.EOF:
		err = errMissingPostContent
	case err != nil:
		err = errors.New("invalid post data")
	case params.Title == "":
		err = errEmptyTitle
	default:
		user := ctx.Value(auth.UserData).(*database.User)
		post, err = s.q.CreatePost(ctx, database.CreatePostParams{
			Title:   params.Title,
			Content: sql.NullString{Valid: true, String: params.Content},
			UserID:  user.ID,
		})
	}

	switch {
	case err != nil:
		internal.RespondWithError(w, http.StatusBadRequest, err.Error())
	default:
		internal.RespondWithJSON(w, http.StatusOK, databasePostToPost(&post))
	}
}

func (s *PostsStorage) GetPost(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postId")
	id, err := strconv.Atoi(postId)

	var post database.Post

	switch {
	case err != nil:
		// do nothing
	default:
		post, err = s.q.GetPost(r.Context(), int32(id))
		if err != nil {
			err = errNotFound
		}
	}

	switch v := err.(type) {
	case *strconv.NumError:
		_ = v
		internal.RespondWithError(w, http.StatusBadRequest, internal.ErrInvalidPostId)
	case error:
		internal.RespondWithError(w, http.StatusBadRequest, err)
	default:
		internal.RespondWithJSON(w, http.StatusOK, databasePostToPost(&post))
	}
}

func (s *PostsStorage) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	// TODO: Offset and Limit via URL query
	// TODO: filter by date and tags
	posts, err := s.q.GetPosts(r.Context(), database.GetPostsParams{
		Offset: 0,
		Limit:  10,
	})

	switch v := err.(type) {
	case *strconv.NumError:
		_ = v
		internal.RespondWithError(w, http.StatusBadRequest, internal.ErrInvalidPostId)
	case error:
		internal.RespondWithError(w, http.StatusBadRequest, "unable get posts")
	default:
		internal.RespondWithJSON(w, http.StatusOK, databasePostsToPosts(&posts))
	}
}

func (s *PostsStorage) UpdatePost(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postId")
	id, err := strconv.Atoi(postId)

	if err != nil {
		internal.RespondWithError(w, http.StatusBadRequest, internal.ErrInvalidPostId)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := changePostRequestParams{}
	err = decoder.Decode(&params)

	var post database.Post
	switch {
	case err == io.EOF:
		err = errMissingPostContent
	case params.Title == "":
		err = errEmptyTitle
	case err == nil:
		userData := ctx.Value(auth.UserData).(*database.User)
		post, err = s.q.UpdatePost(ctx, database.UpdatePostParams{
			ID:      int32(id),
			UserID:  userData.ID,
			Title:   params.Title,
			Content: sql.NullString{String: params.Content, Valid: true},
		})

		if err != nil {
			err = errNotFound
		}
	}

	switch {
	case err != nil:
		internal.RespondWithError(w, http.StatusBadRequest, err.Error())
	default:
		internal.RespondWithJSON(w, http.StatusOK, databasePostToPost(&post))
	}
}

func (s *PostsStorage) DeletePost(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postId")
	id, err := strconv.Atoi(postId)

	switch {
	case err == nil:
		userData := ctx.Value(auth.UserData).(*database.User)
		err = s.q.DeletePost(ctx, database.DeletePostParams{
			ID:     int32(id),
			UserID: userData.ID,
		})
	}

	switch v := err.(type) {
	case *strconv.NumError:
		_ = v
		internal.RespondWithError(w, http.StatusBadRequest, internal.ErrInvalidPostId)
	case error:
		internal.RespondWithError(w, http.StatusBadRequest, "unable to delete post")
	default:
		internal.RespondWithJSON(w, http.StatusOK, struct{}{})
	}
}

func SearchContent(ctx context.Context, w http.ResponseWriter, r *http.Request) {

}
