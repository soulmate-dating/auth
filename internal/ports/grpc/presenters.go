package grpc

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/soulmate-dating/auth/internal/domain"
	"google.golang.org/grpc/codes"
)

var ErrMissingArgument = errors.New("required argument is missing")

func TokenSuccessResponse(p *domain.Token) *TokenResponse {
	return &TokenResponse{
		Id:           p.Id.String(),
		AccessToken:  p.AccessToken,
		RefreshToken: p.RefreshToken,
	}
}

func GetErrorCode(err error) codes.Code {
	switch {
	case errors.As(err, &validator.ValidationErrors{}) || errors.Is(err, domain.ErrFailedToParseClaims):
		return codes.InvalidArgument
	case errors.Is(err, domain.ErrUserNotFound):
		return codes.NotFound
	case errors.Is(err, domain.ErrAlreadyExists):
		return codes.AlreadyExists
	case errors.Is(err, domain.ErrInvalidToken) ||
		errors.Is(err, domain.ErrWrongPassword) ||
		errors.Is(err, domain.ErrExpiredToken):
		return codes.Unauthenticated
	}
	return codes.Internal
}
