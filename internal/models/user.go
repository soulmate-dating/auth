package models

type User struct {
	ID       string `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type UserID struct {
	ID string `db:"id"`
}
