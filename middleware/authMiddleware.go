package middleware

import (
	"context"
	"db-service/internal"
	"db-service/internal/database"
	"fmt"
	"net/http"
)

type AuthService interface {
	AuthorizeUser(ctx context.Context, accessToken string) (*database.User, error)
}

type Middleware interface {
	HandlerFunc(handler middlewareHandler) http.HandlerFunc
}

type middlewareHandler func(w http.ResponseWriter, r *http.Request, ctx context.Context)

type authMiddleware struct {
	Auth AuthService
}

func Auth(authService AuthService) Middleware {
	return &authMiddleware{
		Auth: authService,
	}
}

type contextKey string

const (
	UserData contextKey = "userData"
)

func (m authMiddleware) HandlerFunc(handler middlewareHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := internal.GetAPIKey(r.Header)

		if err != nil {
			internal.RespondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := m.Auth.AuthorizeUser(r.Context(), apiKey)

		if err != nil {
			internal.RespondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}

		ctx := context.WithValue(r.Context(), UserData, user)
		handler(w, r, ctx)
	}
}
