package internal

import (
	"context"
	"net/http"
)

func NotImplemented(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	RespondWithError(w, http.StatusInternalServerError, "not implemented")
}
