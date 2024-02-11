package auth

import "errors"

var (
	errInvalidPasswordLength error = errors.New("invalid password length")
	errUserNotFound          error = errors.New("user not found")
	errCantCreateUser        error = errors.New("user not created")
	errInvalidCredentials    error = errors.New("wrong credentials")
)
