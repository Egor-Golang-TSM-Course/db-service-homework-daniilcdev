package auth

import (
	"db-service/internal"
	"db-service/internal/database"
	"encoding/json"
	"net/http"
	"time"

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
