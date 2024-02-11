package tags

import (
	"context"
	"db-service/handlers/auth"
	"db-service/handlers/posts"
	"db-service/internal"
	"db-service/internal/database"
	"db-service/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const fmtInsertNewTag = `INSERT INTO tags (id, tag, created_at)
VALUES (DEFAULT, UNNEST(ARRAY['%s']), NOW()) ON CONFLICT DO NOTHING
;`

func (ts *TagsStorage) AddTag(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	postId, err := internal.PostIdFromURLParams(r)
	if err != nil {
		internal.RespondWithError(w, http.StatusBadRequest, internal.ErrInvalidPostId)
		return
	}

	var tagsData models.AddedTags
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&tagsData); err != nil || len(tagsData.NewTags) == 0 {
		internal.RespondWithError(w, http.StatusBadRequest, "insufficient tags count: at least 1 expected")
		return
	}

	userData := ctx.Value(auth.UserData).(*database.User)

	post, _ := ts.q.GetPost(ctx, postId)
	if post.UserID != userData.ID {
		internal.RespondWithError(w, http.StatusBadRequest, "invalid ownership")
		return
	}

	tags := internal.Distinct(append(post.Tags, tagsData.NewTags...))

	if len(tags) == len(post.Tags) {
		internal.RespondWithJSON(w, http.StatusAlreadyReported, posts.DatabasePostToPost(&post))
		return
	}

	post, err = ts.q.UpdatePostTags(ctx, database.UpdatePostTagsParams{
		ID:     post.ID,
		UserID: userData.ID,
		Tags:   tags,
	})

	if err != nil {
		internal.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	internal.RespondWithJSON(w, http.StatusOK, posts.DatabasePostToPost(&post))

	timeout, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	// TODO: research how to implement with `sqlc`
	sql := fmt.Sprintf(fmtInsertNewTag, strings.Join(tagsData.NewTags, "','"))
	if err := ts.q.ExecRaw(timeout, sql); err != nil {
		log.Println("failed to update shared Tags collection")
	}
}

func (ts *TagsStorage) GetTags(w http.ResponseWriter, r *http.Request) {
	const batchLimit = 10
	if tags, err := ts.q.AllTags(r.Context(), batchLimit); err != nil {
		internal.RespondWithError(w, http.StatusBadRequest, err)
	} else {
		internal.RespondWithJSON(w, http.StatusOK, databaseTagsToTags(&tags))
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
