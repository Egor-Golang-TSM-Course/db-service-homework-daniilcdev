package auth

import (
	"db-service/internal"
	"db-service/internal/database"
	"db-service/models"

	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (us *UserService) Register(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	params := createUserRequestParams{}
	err := decoder.Decode(&params)

	if err != nil {
		internal.RespondWithError(w, http.StatusBadRequest, internal.ErrInvalidJson)
		return
	}

	pwdBytes := []byte(params.Password)
	if len(pwdBytes) > 72 || len(pwdBytes) == 0 {
		internal.RespondWithError(w, http.StatusBadRequest, errInvalidPasswordLength)
		return
	}

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)

	if err != nil {
		internal.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
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
		internal.RespondWithError(w, http.StatusBadRequest, errCantCreateUser)
		return
	}

	internal.RespondWithJSON(w, http.StatusCreated, databaseUserToUser(user))
}

func (us *UserService) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	params := createUserRequestParams{}
	err := decoder.Decode(&params)

	if err != nil {
		internal.RespondWithError(w, http.StatusBadRequest, internal.ErrInvalidJson)
		return
	}

	pwdBytes := []byte(params.Password)
	if len(pwdBytes) > 72 || len(pwdBytes) == 0 {
		internal.RespondWithError(w, http.StatusBadRequest, errInvalidPasswordLength)
		return
	}

	user, err := us.db.UserByEmail(r.Context(), params.Email)

	if err != nil {
		internal.RespondWithError(w, http.StatusNotFound, errUserNotFound)
		return
	}

	if bcrypt.CompareHashAndPassword(user.PwdHash, pwdBytes) != nil {
		internal.RespondWithJSON(w, http.StatusBadRequest, errInvalidCredentials)
		return
	}

	userData := databaseUserToAuthorizedUser(user)
	internal.RespondWithJSON(w, http.StatusOK, userData)
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
