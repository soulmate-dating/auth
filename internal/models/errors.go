package models

import (
	"errors"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrAlreadyExists       = errors.New("user with email already exists")
	ErrWrongPassword       = errors.New("passwords do not match")
	ErrFailedToParseClaims = errors.New("failed to parse jwt claims")
	ErrInvalidToken        = errors.New("invalid token")
)
