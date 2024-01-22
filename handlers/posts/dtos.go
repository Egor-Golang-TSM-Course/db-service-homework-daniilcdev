package posts

import (
	"db-service/internal/database"

	"github.com/google/uuid"
)

type Post struct {
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	AuthorId uuid.UUID `json:"author_id"`
	Tags     []string  `json:"tags"`
}

func databasePostToPost(post *database.Post) Post {
	return Post{
		Title:    post.Title,
		Content:  post.Content.String,
		AuthorId: post.UserID,
	}
}

func databasePostsToPosts(posts *[]database.Post) []Post {
	r := make([]Post, 0, len(*posts))

	for _, post := range *posts {
		r = append(r, Post{
			Title:    post.Title,
			Content:  post.Content.String,
			AuthorId: post.UserID,
		})
	}
	return r
}
