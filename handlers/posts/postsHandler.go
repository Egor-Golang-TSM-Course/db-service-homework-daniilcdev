package posts

import (
	"context"
	"database/sql"
	"db-service/handlers/auth"
	"db-service/internal"
	"db-service/internal/database"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type PostsStorage struct {
	q *database.Queries
}

func NewStorage(q *database.Queries) *PostsStorage {
	return &PostsStorage{
		q: q,
	}
}

func (s *PostsStorage) CreatePost(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	decoder := json.NewDecoder(r.Body)

	params := createPostRequestParams{}
	err := decoder.Decode(&params)
	var post database.Post

	switch {
	case err != nil:
		internal.RespondWithError(w, 400, "invalid post data")
	case params.Title == "":
		internal.RespondWithError(w, 400, "post 'title' is empty")
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
		internal.RespondWithError(w, 500, err.Error())
	default:
		internal.RespondWithJSON(w, 200, databasePostToPost(&post))
	}
}

func (s *PostsStorage) GetPost(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	postId := chi.URLParam(r, "postId")
	id, err := strconv.Atoi(postId)

	var post database.Post

	switch {
	case err != nil:
		// do nothing
	default:
		post, err = s.q.GetPost(ctx, int32(id))
	}

	switch v := err.(type) {
	case *strconv.NumError:
		_ = v
		internal.RespondWithError(w, 400, "invalid 'postId'")
	case error:
		internal.RespondWithError(w, 400, "post not found")
	default:
		internal.RespondWithJSON(w, 200, databasePostToPost(&post))
	}
}

func (s *PostsStorage) GetAllPosts(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	// TODO: Offset and Limit via URL query
	// TODO: filter by date and tags
	posts, err := s.q.GetPosts(ctx, database.GetPostsParams{
		Offset: 0,
		Limit:  10,
	})

	switch v := err.(type) {
	case *strconv.NumError:
		_ = v
		internal.RespondWithError(w, 400, "invalid 'postId'")
	case error:
		internal.RespondWithError(w, 400, "unable get posts")
	default:
		internal.RespondWithJSON(w, 200, databasePostsToPosts(&posts))
	}
}

func (s *PostsStorage) UpdatePost(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	postId := chi.URLParam(r, "postId")
	id, err := strconv.Atoi(postId)

	if err != nil {
		internal.RespondWithError(w, 400, "invalid 'postId'")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := changePostRequestParams{}
	err = decoder.Decode(&params)

	var post database.Post
	switch {
	case params.Title == "":
		err = errors.New("restriction violation: title can't be empty")
	case err == nil:
		userData := ctx.Value(auth.UserData).(*database.User)
		post, err = s.q.UpdatePost(ctx, database.UpdatePostParams{
			ID:      int32(id),
			UserID:  userData.ID,
			Title:   params.Title,
			Content: sql.NullString{String: params.Content, Valid: true},
		})
	}

	switch {
	case err != nil:
		internal.RespondWithError(w, 400, "unable to update post")
	default:
		internal.RespondWithJSON(w, 200, databasePostToPost(&post))
	}
}

func (s *PostsStorage) DeletePost(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	postId := chi.URLParam(r, "postId")
	id, err := strconv.Atoi(postId)

	var post database.Post

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
		internal.RespondWithError(w, 400, "invalid 'postId'")
	case error:
		internal.RespondWithError(w, 400, "unable to delete post")
	default:
		internal.RespondWithJSON(w, 200, databasePostToPost(&post))
	}
}

func SearchContent(w http.ResponseWriter, r *http.Request, ctx context.Context) {

}
