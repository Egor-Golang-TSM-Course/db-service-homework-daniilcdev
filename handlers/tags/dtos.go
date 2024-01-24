package tags

import (
	"db-service/internal/database"
	"time"
)

type Tag struct {
	Id        int32     `json:"id"`
	Tag       string    `json:"tag"`
	CreatedAt time.Time `json:"created_at"`
}

type addTagToPostDto struct {
	NewTags []string `json:"newTags"`
}

func databaseTagToTag(tag *database.Tag) Tag {
	return Tag{
		Id:        tag.ID,
		Tag:       tag.Tag,
		CreatedAt: tag.CreatedAt,
	}
}

func databaseTagsToTags(tags *[]database.Tag) []Tag {
	out := make([]Tag, 0, len(*tags))

	for _, tag := range *tags {
		out = append(out, databaseTagToTag(&tag))
	}

	return out
}
