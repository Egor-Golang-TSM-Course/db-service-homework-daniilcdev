package tags

import (
	"context"
	"db-service/handlers"
	"db-service/handlers/auth"
	"db-service/handlers/posts"
	"db-service/internal"
	"db-service/internal/database"
	"db-service/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const fmtInsertNewTag = `INSERT INTO tags (id, tag, created_at)
VALUES (DEFAULT, UNNEST(ARRAY['%s']), NOW()) ON CONFLICT DO NOTHING
;`

func (ts *TagsStorage) AddTag(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	postId, err := internal.PostIdFromURLParams(r)
	if err != nil {
		return &handlers.HttpError{Code: http.StatusBadRequest, Err: internal.ErrInvalidPostId}
	}

	var tagsData models.AddedTags
	if err = json.NewDecoder(r.Body).Decode(&tagsData); err != nil || len(tagsData.NewTags) == 0 {
		return &handlers.HttpError{Code: http.StatusBadRequest, Err: errors.New("insufficient tags count: at least 1 expected")}
	}

	userData := ctx.Value(auth.UserData).(*database.User)

	post, err := ts.q.GetPost(ctx, postId)

	if err != nil {
		return &handlers.HttpError{Code: http.StatusBadRequest, Err: errors.New("post not found")}
	}

	if post.UserID != userData.ID {
		return &handlers.HttpError{Code: http.StatusBadRequest, Err: errors.New("invalid ownership")}
	}

	tags := internal.Distinct(append(post.Tags, tagsData.NewTags...))

	if len(tags) == len(post.Tags) {
		fmt.Fprint(w, posts.DatabasePostToPost(&post))
		internal.SetHeaders(w, http.StatusAlreadyReported)
		return nil
	}

	post, err = ts.q.UpdatePostTags(ctx, database.UpdatePostTagsParams{
		ID:     post.ID,
		UserID: userData.ID,
		Tags:   tags,
	})

	if err != nil {
		return &handlers.HttpError{Code: http.StatusInternalServerError, Err: err}
	}

	fmt.Fprint(w, posts.DatabasePostToPost(&post))
	internal.SetHeaders(w, http.StatusOK)

	timeout, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	// TODO: research how to implement with `sqlc`
	sql := fmt.Sprintf(fmtInsertNewTag, strings.Join(tagsData.NewTags, "','"))
	if err := ts.q.ExecRaw(timeout, sql); err != nil {
		log.Println("failed to update shared Tags collection")
	}

	return nil
}

func (ts *TagsStorage) GetTags(w http.ResponseWriter, r *http.Request) error {
	const batchLimit = 10
	if tags, err := ts.q.AllTags(r.Context(), batchLimit); err != nil {
		return &handlers.HttpError{Code: http.StatusBadRequest, Err: err}
	} else {
		fmt.Fprint(w, databaseTagsToTags(&tags))
		internal.SetHeaders(w, http.StatusOK)
		return nil
	}
}

func databaseTagToTag(tag *database.Tag) models.Tag {
	return models.Tag{
		Id:        tag.ID,
		Tag:       tag.Tag,
		CreatedAt: tag.CreatedAt,
	}
}

func databaseTagsToTags(tags *[]database.Tag) []models.Tag {
	out := make([]models.Tag, 0, len(*tags))

	for _, tag := range *tags {
		out = append(out, databaseTagToTag(&tag))
	}

	return out
}
