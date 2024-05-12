package domain

import "github.com/google/uuid"

type Token struct {
	Id           uuid.UUID `json:"id" validate:"required"`
	AccessToken  string    `json:"access_token" validate:"required,jwt"`
	RefreshToken string    `json:"refresh_token" validate:"required,jwt"`
}
