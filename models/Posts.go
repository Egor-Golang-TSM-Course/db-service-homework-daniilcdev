package models

import (
	"github.com/google/uuid"
)

type Post struct {
	Id       int32     `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	AuthorId uuid.UUID `json:"author_id"`
	Tags     []string  `json:"tags"`
}
