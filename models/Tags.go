package models

import (
	"time"
)

type Tag struct {
	Id        int32     `json:"id"`
	Tag       string    `json:"tag"`
	CreatedAt time.Time `json:"created_at"`
}

type AddedTags struct {
	NewTags []string `json:"newTags"`
}
