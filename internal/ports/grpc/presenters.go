package grpc

import (
	"errors"
	"github.com/TobbyMax/validator"
	"github.com/soulmate-dating/auth/internal/app"
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
	case errors.As(err, &validator.ValidationErrors{}):
		return codes.InvalidArgument
	case errors.Is(err, app.ErrForbidden):
		return codes.PermissionDenied
	}
	return codes.Internal
}
