package auth

import (
	"context"
	"db-service/internal"
	"db-service/internal/database"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type AuthService interface {
	AuthorizeUser(ctx context.Context, accessToken string) (*database.User, error)
}

type Middleware interface {
	HandlerFunc(handler middlewareHandler) http.HandlerFunc
}

type middlewareHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request)

type authMiddleware struct {
	Auth AuthService
}

func NewMiddleware(authService AuthService) Middleware {
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
		accessToken, err := getAccessToken(r.Header)

		if err != nil {
			internal.RespondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		err = VerifyAccessToken(accessToken)
		if err != nil {
			internal.RespondWithError(w, http.StatusUnauthorized, "invalid access token")
			return
		}

		user, err := m.Auth.AuthorizeUser(r.Context(), accessToken)

		if err != nil {
			internal.RespondWithError(w, 400, fmt.Sprintf("unauthorized access %v", err))
			return
		}

		ctx := context.WithValue(r.Context(), UserData, user)
		handler(ctx, w, r)
	}
}

func (us *UserService) AuthorizeUser(ctx context.Context, accessToken string) (*database.User, error) {
	claims := ParseAccessToken(accessToken)
	user, err := us.db.UserById(ctx, claims.Id)
	return &user, err
}

func getAccessToken(headers http.Header) (string, error) {

	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentication info found")
	}

	vals := strings.Split(val, " ")

	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}

	if vals[0] != "Bearer" {
		return "", errors.New("malformed first part of auth header")
	}

	return vals[1], nil
}
