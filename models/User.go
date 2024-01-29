package models

import (
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
