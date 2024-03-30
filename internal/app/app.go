package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/soulmate-dating/auth/internal/adapters/jwt"
	"github.com/soulmate-dating/auth/internal/adapters/postgres"
	"github.com/soulmate-dating/auth/internal/hash"
	"github.com/soulmate-dating/auth/internal/models"
)

var (
	ErrForbidden = fmt.Errorf("forbidden")
)

type App interface {
	SignUp(ctx context.Context, email string, password string) (*models.Token, error)
	Login(ctx context.Context, email string, password string) (*models.Token, error)
	Refresh(ctx context.Context, token string) (*models.Token, error)
	Logout(ctx context.Context, token string) (string, error)
	Validate(ctx context.Context, token string) (string, error)
}

type Repository interface {
	CreateUser(ctx context.Context, p *models.User) (uuid.UUID, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	//UpdateLoginStatus(ctx context.Context, uid string, b bool) error
}

type Application struct {
	repository Repository
	jwtWrapper jwt.Wrapper
}

func (a *Application) Validate(ctx context.Context, token string) (string, error) {
	claims, err := a.jwtWrapper.ValidateToken(token)
	if err != nil {
		return "", err
	}
	return claims.Id.String(), nil
}

func (a *Application) Logout(ctx context.Context, token string) (string, error) {
	claims, err := a.jwtWrapper.ValidateToken(token)
	if err != nil {
		return "", err
	}
	return claims.Id.String(), nil
}

func (a *Application) SignUp(ctx context.Context, email string, password string) (*models.Token, error) {
	user := &models.User{
		ID:       models.NewUUID(),
		Email:    email,
		Password: hash.HashPassword(password),
	}

	id, err := a.repository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	accessToken, err := a.jwtWrapper.GenerateAccessToken(user)
	if err != nil {
		return nil, errors.New("error generating token")
	}

	refreshToken, err := a.jwtWrapper.GenerateRefreshToken(user)
	if err != nil {
		return nil, errors.New("error generating token")
	}

	return &models.Token{Id: id, AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (a *Application) Login(ctx context.Context, email string, password string) (*models.Token, error) {
	user, err := a.repository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	match := hash.CheckPasswordHash(password, user.Password)
	if !match {
		return nil, ErrForbidden
	}

	accessToken, err := a.jwtWrapper.GenerateAccessToken(user)
	if err != nil {
		return nil, errors.New("error generating token")
	}

	refreshToken, err := a.jwtWrapper.GenerateRefreshToken(user)
	if err != nil {
		return nil, errors.New("error generating token")
	}

	return &models.Token{Id: user.ID, AccessToken: accessToken, RefreshToken: refreshToken}, nil

}

func (a *Application) Refresh(ctx context.Context, token string) (*models.Token, error) {
	claims, err := a.jwtWrapper.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	user, err := a.repository.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return nil, err
	}

	accessToken, err := a.jwtWrapper.GenerateAccessToken(user)
	if err != nil {
		return nil, errors.New("error generating token")
	}
	refreshToken, err := a.jwtWrapper.GenerateRefreshToken(user)
	if err != nil {
		return nil, errors.New("error generating token")
	}

	return &models.Token{Id: user.ID, AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func NewApp(conn *pgxpool.Pool, jwt *jwt.Wrapper) App {
	repo := postgres.NewRepo(conn)
	return &Application{repository: repo, jwtWrapper: *jwt}
}
