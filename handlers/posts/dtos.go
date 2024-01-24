package posts

import (
	"db-service/internal/database"

	"github.com/google/uuid"
)

type Post struct {
	Id       int32     `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	AuthorId uuid.UUID `json:"author_id"`
	Tags     []string  `json:"tags"`
}

func DatabasePostToPost(post *database.Post) Post {
	return Post{
		Id:       post.ID,
		Title:    post.Title,
		Content:  post.Content.String,
		AuthorId: post.UserID,
		Tags:     post.Tags,
	}
}

func databasePostsToPosts(posts *[]database.Post) []Post {
	r := make([]Post, 0, len(*posts))

	for _, post := range *posts {
		r = append(r, DatabasePostToPost(&post))
	}
	return r
}
