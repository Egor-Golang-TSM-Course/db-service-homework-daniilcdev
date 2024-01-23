package auth

import (
	"db-service/internal/database"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type User struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type AuthorizedUser struct {
	User
	Token string `json:"token"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:    dbUser.ID,
		Name:  dbUser.Name,
		Email: dbUser.Email,
	}
}

func databaseUserToAuthorizedUser(dbUser database.User) AuthorizedUser {
	user := AuthorizedUser{}
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
