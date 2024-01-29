package posts

import (
	"context"
	"database/sql"
	"db-service/handlers/auth"
	"db-service/internal"
	"db-service/internal/database"
	"db-service/models"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
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
		if err != nil {
			err = errNotFound
		}
	}

	switch {
	case err == errNotFound:
		internal.RespondWithError(w, http.StatusNotFound, err)
	case err != nil:
		internal.RespondWithError(w, http.StatusBadRequest, err)
	default:
		internal.RespondWithJSON(w, http.StatusOK, DatabasePostToPost(&post))
	}
}

func (s *PostsStorage) GetPost(w http.ResponseWriter, r *http.Request) {
	postId, err := internal.PostIdFromURLParams(r)

	if err != nil {
		internal.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	post, err := s.q.GetPost(r.Context(), postId)

	if err != nil {
		err = errNotFound
	}

	switch {
	case err == errNotFound:
		internal.RespondWithError(w, http.StatusNotFound, internal.ErrInvalidPostId)
	case err != nil:
		internal.RespondWithError(w, http.StatusBadRequest, err)
	default:
		internal.RespondWithJSON(w, http.StatusOK, DatabasePostToPost(&post))
	}
}

func (s *PostsStorage) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	// TODO: Offset and Limit via URL query
	// TODO: filter by date and tags
	const offset = 0
	const limit = 10

	posts, err := s.q.GetPosts(r.Context(), database.GetPostsParams{
		Offset: offset,
		Limit:  limit,
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
	postId, err := internal.PostIdFromURLParams(r)

	if err != nil {
		internal.RespondWithError(w, http.StatusBadRequest, err)
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
			ID:      postId,
			UserID:  userData.ID,
			Title:   params.Title,
			Content: sql.NullString{String: params.Content, Valid: true},
		})

		if err != nil {
			err = errNotFound
		}
	}

	switch {
	case err == errNotFound:
		internal.RespondWithError(w, http.StatusNotFound, err)
	case err != nil:
		internal.RespondWithError(w, http.StatusBadRequest, err)
	default:
		internal.RespondWithJSON(w, http.StatusOK, DatabasePostToPost(&post))
	}
}

func (s *PostsStorage) DeletePost(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	postId, err := internal.PostIdFromURLParams(r)

	switch {
	case err == nil:
		userData := ctx.Value(auth.UserData).(*database.User)
		err = s.q.DeletePost(ctx, database.DeletePostParams{
			ID:     postId,
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

func DatabasePostToPost(post *database.Post) models.Post {
	return models.Post{
		Id:       post.ID,
		Title:    post.Title,
		Content:  post.Content.String,
		AuthorId: post.UserID,
		Tags:     post.Tags,
	}
}

func databasePostsToPosts(posts *[]database.Post) []models.Post {
	r := make([]models.Post, 0, len(*posts))

	for _, post := range *posts {
		r = append(r, DatabasePostToPost(&post))
	}
	return r
}
