package handlers

import (
	"context"
	"db-service/internal"
	"net/http"
)

type HttpError struct {
	Code int
	Err  error
}

func (e *HttpError) Error() string {
	return e.Err.Error()
}

type Wrapper struct {
}

func Wrap(base func(w http.ResponseWriter, r *http.Request) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := base(w, r)
		if err != nil {
			switch e := err.(type) {
			case *HttpError:
				internal.RespondWithError(w, e.Code, err)
			default:
				internal.RespondWithError(w, http.StatusBadRequest, err)
			}
		}
	}
}

func WrapCtx(base func(ctx context.Context, w http.ResponseWriter, r *http.Request) error) func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		err := base(ctx, w, r)
		if err != nil {
			switch e := err.(type) {
			case *HttpError:
				internal.RespondWithError(w, e.Code, err)
			default:
				internal.RespondWithError(w, http.StatusBadRequest, err)
			}

		}
	}
}
