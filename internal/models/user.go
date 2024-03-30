package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `db:"id"`
	Email    string    `db:"email"`
	Password string    `db:"password"`
}

type UserID struct {
	ID uuid.UUID `db:"id"`
}
