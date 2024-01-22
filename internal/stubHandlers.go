package internal

import (
	"context"
	"net/http"
)

func NotImplemented(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	RespondWithError(w, 500, "not implemented")
}
