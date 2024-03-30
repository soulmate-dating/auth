package models

import "github.com/google/uuid"

type Token struct {
	Id           uuid.UUID
	AccessToken  string
	RefreshToken string
}
