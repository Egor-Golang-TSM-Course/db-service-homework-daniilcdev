package auth

import (
	"db-service/internal"
	"db-service/internal/database"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db *database.Queries
}

func NewService() *UserService {
	return &UserService{}
}

func (us *UserService) WithDb(db *database.Queries) *UserService {
	us.db = db
	return us
}

func (us *UserService) Register(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	params := createUserRequestParams{}
	err := decoder.Decode(&params)

	if err != nil {
		internal.RespondWithError(w, 400, fmt.Sprint("Error parsing JSON", err))
		return
	}

	pwdBytes := []byte(params.Password)
	if len(pwdBytes) > 72 {
		internal.RespondWithError(w, 400, "password is too long")
		return
	}

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)

	if err != nil {
		internal.RespondWithError(w, 500, err.Error())
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
		internal.RespondWithError(w, 400, fmt.Sprint("user not created:", err))
		return
	}

	internal.RespondWithJSON(w, 201, databaseUserToUser(user))
}

func (us *UserService) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	params := createUserRequestParams{}
	err := decoder.Decode(&params)

	if err != nil {
		internal.RespondWithError(w, 400, fmt.Sprint("Error parsing JSON", err))
		return
	}

	pwdBytes := []byte(params.Password)
	if len(pwdBytes) > 72 {
		internal.RespondWithError(w, 400, "password is too long")
		return
	}

	user, err := us.db.UserByEmail(r.Context(), params.Email)

	if err != nil {
		internal.RespondWithError(w, 404, fmt.Sprint("user not found:", err))
		return
	}

	if bcrypt.CompareHashAndPassword(user.PwdHash, pwdBytes) != nil {
		internal.RespondWithJSON(w, 400, "wrong credentials")
		return
	}

	userData := databaseUserToAuthorizedUser(user)
	internal.RespondWithJSON(w, 200, userData)
}
