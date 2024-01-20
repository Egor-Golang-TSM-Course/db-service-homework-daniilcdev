package auth

import (
	"github.com/google/uuid"
)

type User struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Token string    `json:"token"`
}
