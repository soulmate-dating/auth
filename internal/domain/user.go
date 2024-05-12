package domain

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `db:"id" validate:"required"`
	Email    string    `db:"email" validate:"required,email"`
	Password string    `db:"password" validate:"required,min=8"`
}

type LoginCredentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserID struct {
	ID uuid.UUID `db:"id"`
}
