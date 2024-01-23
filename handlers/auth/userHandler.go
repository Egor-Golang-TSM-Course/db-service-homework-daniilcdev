package auth

import (
	"context"

	"db-service/internal"
	"db-service/internal/database"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
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

	userData := databaseUserToUser(user)
	accessToken, _ := NewAccessToken(Claims{
		Id:  user.ID.String(),
		Key: user.PwdHash,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	})
	userData.Token = accessToken

	internal.RespondWithJSON(w, 200, userData)
}

func (us *UserService) AuthorizeUser(ctx context.Context, accessToken string) (*database.User, error) {
	claims := ParseAccessToken(accessToken)
	user, err := us.db.UserByAuthToken(ctx, claims.Key)
	return &user, err
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:    dbUser.ID,
		Name:  dbUser.Name,
		Email: dbUser.Email,
	}
}
