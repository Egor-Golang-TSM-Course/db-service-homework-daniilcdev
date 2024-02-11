package internal

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func PostIdFromURLParams(r *http.Request) (int32, error) {
	postId := chi.URLParam(r, "postId")
	id, err := strconv.Atoi(postId)
	if err != nil {
		return -1, ErrInvalidPostId
	}
	return int32(id), nil

}
