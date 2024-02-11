package models

import (
	"time"

	"github.com/google/uuid"
)

type Comments struct {
	Comments []Comment `json:"comments"`
}

type Comment struct {
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	Author    Author    `json:"author"`
}

type Author struct {
	Id uuid.UUID `json:"id"`
}
