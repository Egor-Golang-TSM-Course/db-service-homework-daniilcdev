package userService

import (
	"context"
	"crypto"
	"db-service/internal"
	"db-service/internal/database"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
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

	user, err := us.db.CreateUser(r.Context(), database.CreateUserParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Name:        params.Name,
		Email:       params.Email,
		AccessToken: pwdHashString(params.Email + ":" + params.Password),
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

	user, err := us.db.AuthorizeUser(r.Context(), pwdHashString(params.Email+":"+params.Password))

	if err != nil {
		internal.RespondWithError(w, 400, fmt.Sprint("user not created:", err))
		return
	}

	internal.RespondWithJSON(w, 200, databaseUserToUser(user))
}

func (us *UserService) AuthorizeUser(ctx context.Context, accessToken string) (*database.User, error) {
	user, err := us.db.AuthorizeUser(ctx, accessToken)
	return &user, err
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:          dbUser.ID,
		Name:        dbUser.Name,
		Email:       dbUser.Email,
		AccessToken: dbUser.AccessToken,
	}
}

func pwdHashString(raw string) string {
	var bytes []byte = []byte(raw)
	sha256 := crypto.SHA256.New()
	sha256.Write(bytes)
	return base64.URLEncoding.EncodeToString(sha256.Sum(nil))
}
