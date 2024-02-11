package posts

import (
	"context"
	"database/sql"
	"db-service/handlers"
	"db-service/handlers/auth"
	"db-service/internal"
	"db-service/internal/database"
	"db-service/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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

func (s *PostsStorage) CreatePost(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	params := createPostRequestParams{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return &handlers.HttpError{Code: http.StatusBadRequest, Err: errors.New("invalid post data")}
	}

	if params.Title == "" {
		return &handlers.HttpError{Code: http.StatusBadRequest, Err: errEmptyTitle}
	}

	user := ctx.Value(auth.UserData).(*database.User)
	post, err := s.q.CreatePost(ctx, database.CreatePostParams{
		Title:   params.Title,
		Content: sql.NullString{Valid: true, String: params.Content},
		UserID:  user.ID,
	})

	if err != nil {
		return &handlers.HttpError{Code: http.StatusNotFound, Err: errNotFound}
	}

	fmt.Fprint(w, DatabasePostToPost(&post))
	internal.SetHeaders(w, http.StatusCreated)
	return nil
}

func (s *PostsStorage) GetPost(w http.ResponseWriter, r *http.Request) error {
	postId, err := internal.PostIdFromURLParams(r)

	if err != nil {
		return &handlers.HttpError{Code: http.StatusBadRequest, Err: err}
	}

	post, err := s.q.GetPost(r.Context(), postId)

	if err != nil {
		return &handlers.HttpError{Code: http.StatusNotFound, Err: errNotFound}
	}

	fmt.Fprint(w, DatabasePostToPost(&post))
	internal.SetHeaders(w, http.StatusOK)
	return nil
}

func (s *PostsStorage) GetAllPosts(w http.ResponseWriter, r *http.Request) error {
	// TODO: Offset and Limit via URL query
	// TODO: filter by date and tags
	const offset = 0
	const limit = 10

	posts, err := s.q.GetPosts(r.Context(), database.GetPostsParams{
		Offset: offset,
		Limit:  limit,
	})

	if err != nil {
		return &handlers.HttpError{Code: http.StatusBadRequest, Err: err}
	}

	fmt.Fprint(w, databasePostsToPosts(&posts))
	internal.SetHeaders(w, http.StatusOK)
	return nil
}

func (s *PostsStorage) UpdatePost(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	postId, err := internal.PostIdFromURLParams(r)

	if err != nil {
		return &handlers.HttpError{Code: http.StatusBadRequest, Err: err}
	}

	params := changePostRequestParams{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return &handlers.HttpError{Code: http.StatusBadRequest, Err: errMissingPostContent}
	}

	if params.Title == "" {
		return &handlers.HttpError{Code: http.StatusBadRequest, Err: errEmptyTitle}
	}

	userData := ctx.Value(auth.UserData).(*database.User)
	post, err := s.q.UpdatePost(ctx, database.UpdatePostParams{
		ID:      postId,
		UserID:  userData.ID,
		Title:   params.Title,
		Content: sql.NullString{String: params.Content, Valid: true},
	})

	if err != nil {
		return &handlers.HttpError{Code: http.StatusNotFound, Err: errNotFound}
	}

	fmt.Fprint(w, DatabasePostToPost(&post))
	internal.SetHeaders(w, http.StatusOK)
	return nil
}

func (s *PostsStorage) DeletePost(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	postId, err := internal.PostIdFromURLParams(r)

	if err != nil {
		return &handlers.HttpError{Code: http.StatusBadRequest, Err: internal.ErrInvalidPostId}
	}

	userData := ctx.Value(auth.UserData).(*database.User)
	if err = s.q.DeletePost(ctx, database.DeletePostParams{
		ID:     postId,
		UserID: userData.ID,
	}); err != nil {
		return &handlers.HttpError{Code: http.StatusBadRequest, Err: err}
	}

	internal.SetHeaders(w, http.StatusOK)
	return nil
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
