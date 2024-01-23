package auth

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Claims struct {
	Id uuid.UUID `json:"id"`
	jwt.StandardClaims
}

func NewAccessToken(claims Claims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func NewRefreshToken(claims jwt.StandardClaims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return refreshToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func ParseAccessToken(accessToken string) *Claims {
	parsedAccessToken, _ := jwt.ParseWithClaims(accessToken, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TOKEN_SECRET")), nil
		})

	return parsedAccessToken.Claims.(*Claims)
}

func VerifyAccessToken(accessToken string) error {
	token, err := jwt.Parse(accessToken,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TOKEN_SECRET")), nil
		},
	)

	switch {
	case err != nil:
		return err
	case !token.Valid:
		return fmt.Errorf("invalid token")
	default:
		return nil
	}
}

func ParseRefreshToken(refreshToken string) *jwt.StandardClaims {
	parsedRefreshToken, _ := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TOKEN_SECRET")), nil
		})

	return parsedRefreshToken.Claims.(*jwt.StandardClaims)
}
