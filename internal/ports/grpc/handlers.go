package grpc

import (
	"context"

	"google.golang.org/grpc/status"

	"github.com/soulmate-dating/auth/internal/models"
)

func (s *AuthService) SignUp(ctx context.Context, request *SignUpRequest) (*TokenResponse, error) {
	token, err := s.app.SignUp(ctx, models.LoginCredentials{
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	})
	if err != nil {
		return nil, status.Error(GetErrorCode(err), err.Error())
	}
	return TokenSuccessResponse(token), nil
}

func (s *AuthService) Login(ctx context.Context, request *LoginRequest) (*TokenResponse, error) {
	token, err := s.app.Login(ctx, models.LoginCredentials{
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	})
	if err != nil {
		return nil, status.Error(GetErrorCode(err), err.Error())
	}
	return TokenSuccessResponse(token), nil
}

func (s *AuthService) Logout(ctx context.Context, request *LogoutRequest) (*UserResponse, error) {
	id, err := s.app.Logout(ctx, request.GetAccessToken())
	if err != nil {
		return nil, status.Error(GetErrorCode(err), err.Error())
	}
	return &UserResponse{Id: id}, nil
}

func (s *AuthService) Validate(ctx context.Context, request *ValidateRequest) (*UserResponse, error) {
	id, err := s.app.Validate(ctx, request.GetAccessToken())
	if err != nil {
		return nil, status.Error(GetErrorCode(err), err.Error())
	}
	return &UserResponse{Id: id}, nil
}

func (s *AuthService) Refresh(ctx context.Context, request *RefreshRequest) (*TokenResponse, error) {
	token, err := s.app.Refresh(ctx, request.GetRefreshToken())
	if err != nil {
		return nil, status.Error(GetErrorCode(err), err.Error())
	}
	return TokenSuccessResponse(token), nil
}

func (s *AuthService) mustEmbedUnimplementedAuthServiceServer() {}
