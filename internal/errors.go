package internal

import "errors"

var (
	ErrInvalidPostId = errors.New("invalid 'postId'")
	ErrInvalidJson   = errors.New("invalid JSON format")
)
