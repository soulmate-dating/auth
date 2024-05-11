package grpc

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/soulmate-dating/auth/internal/models"
	"google.golang.org/grpc/codes"
)

var ErrMissingArgument = errors.New("required argument is missing")

func TokenSuccessResponse(p *models.Token) *TokenResponse {
	return &TokenResponse{
		Id:           p.Id.String(),
		AccessToken:  p.AccessToken,
		RefreshToken: p.RefreshToken,
	}
}

func GetErrorCode(err error) codes.Code {
	switch {
	case errors.As(err, &validator.ValidationErrors{}) || errors.Is(err, models.ErrFailedToParseClaims):
		return codes.InvalidArgument
	case errors.Is(err, models.ErrUserNotFound):
		return codes.NotFound
	case errors.Is(err, models.ErrAlreadyExists):
		return codes.AlreadyExists
	case errors.Is(err, models.ErrInvalidToken) ||
		errors.Is(err, models.ErrWrongPassword) ||
		errors.Is(err, models.ErrExpiredToken):
		return codes.Unauthenticated
	}
	return codes.Internal
}
