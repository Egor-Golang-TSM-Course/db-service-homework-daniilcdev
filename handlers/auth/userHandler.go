package auth

import (
	"db-service/handlers"
	"db-service/internal"
	"db-service/internal/database"
	"db-service/models"
	"fmt"

	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (us *UserService) Register(w http.ResponseWriter, r *http.Request) error {
	params := createUserRequestParams{}
	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		return &handlers.HttpError{Code: http.StatusBadRequest, Err: err}
	}

	pwdBytes := []byte(params.Password)
	if len(pwdBytes) > 72 || len(pwdBytes) == 0 {
		return &handlers.HttpError{Code: http.StatusBadRequest, Err: errInvalidPasswordLength}
	}

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)

	if err != nil {
		return &handlers.HttpError{Code: http.StatusInternalServerError, Err: err}
	}

	user, err := us.db.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Email:     params.Email,
		PwdHash:   pwdHash,
	})

	if err != nil {
		return &handlers.HttpError{Code: http.StatusBadRequest, Err: errCantCreateUser}
	}

	fmt.Fprint(w, databaseUserToUser(user))
	internal.SetHeaders(w, http.StatusOK)
	return nil
}

func (us *UserService) Login(w http.ResponseWriter, r *http.Request) error {
	params := createUserRequestParams{}
	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		return &handlers.HttpError{Code: http.StatusBadRequest, Err: internal.ErrInvalidJson}
	}

	pwdBytes := []byte(params.Password)
	if len(pwdBytes) > 72 || len(pwdBytes) == 0 {
		return &handlers.HttpError{Code: http.StatusBadRequest, Err: errInvalidPasswordLength}
	}

	user, err := us.db.UserByEmail(r.Context(), params.Email)

	if err != nil {
		return &handlers.HttpError{Code: http.StatusNotFound, Err: errUserNotFound}
	}

	if bcrypt.CompareHashAndPassword(user.PwdHash, pwdBytes) != nil {
		return &handlers.HttpError{Code: http.StatusBadRequest, Err: errInvalidCredentials}
	}

	fmt.Fprint(w, databaseUserToAuthorizedUser(user))
	internal.SetHeaders(w, http.StatusOK)
	return nil
}

func databaseUserToUser(dbUser database.User) models.User {
	return models.User{
		ID:    dbUser.ID,
		Name:  dbUser.Name,
		Email: dbUser.Email,
	}
}

func databaseUserToAuthorizedUser(dbUser database.User) models.AuthorizedUser {
	user := models.AuthorizedUser{}
	user.ID = dbUser.ID
	user.Name = dbUser.Name
	user.Email = dbUser.Email

	accessToken, _ := NewAccessToken(Claims{
		Id: user.ID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	})

	user.Token = accessToken
	return user
}
