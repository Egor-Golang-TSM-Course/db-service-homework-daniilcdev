package tags

import (
	"context"
	"db-service/handlers/auth"
	"db-service/handlers/posts"
	"db-service/internal"
	"db-service/internal/database"
	"encoding/json"
	"fmt"
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

	var tagData addTagToPostDto
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&tagData); err != nil || len(tagData.NewTags) == 0 {
		internal.RespondWithError(w, http.StatusBadRequest, "insufficient tags count: at least 1 expected")
		return
	}

	userData := ctx.Value(auth.UserData).(*database.User)

	post, _ := ts.q.GetPost(ctx, postId)
	if post.UserID != userData.ID {
		internal.RespondWithError(w, http.StatusBadRequest, "invalid ownership")
		return
	}

	tags := append(post.Tags, tagData.NewTags...)

	{ // remove duplicates
		unique := make(map[string]struct{}, len(post.Tags))

		for _, tag := range tags {
			unique[tag] = struct{}{}
		}

		tags = tags[:0]
		fmt.Println(unique)

		for tag := range unique {
			tags = append(tags, tag)
		}
		fmt.Println(tags)
	}

	if len(tags) == len(post.Tags) {
		internal.RespondWithJSON(w, http.StatusAlreadyReported, posts.DatabasePostToPost(&post))
		return
	}

	post, err = ts.q.UpdatePostTags(ctx, database.UpdatePostTagsParams{
		ID:     post.ID,
		UserID: userData.ID,
		Tags:   tags,
	})

	switch {
	case err != nil:
		internal.RespondWithError(w, http.StatusInternalServerError, err)
	default:
		internal.RespondWithJSON(w, http.StatusOK, posts.DatabasePostToPost(&post))

		timeout, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
		defer cancel()

		// TODO: research how to implement with `sqlc`
		sql := fmt.Sprintf(fmtInsertNewTag, strings.Join(tagData.NewTags, "','"))
		ts.q.ExecRaw(timeout, sql)
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
