package user

import "errors"

var (
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrRoleNotFound          = errors.New("role not found")
	ErrTokenNotFound         = errors.New("token not found")
	ErrInvalidToken          = errors.New("invalid token")
	ErrTokenExpired          = errors.New("token expired")
	ErrInvalidCredentials    = errors.New("invalid credentials")
)
