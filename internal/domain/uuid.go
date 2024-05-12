package domain

import "github.com/google/uuid"

func NewUUID() uuid.UUID {
	return uuid.New()
}
