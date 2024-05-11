package app

import (
	"context"
	"fmt"
	"github.com/google/uuid"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/soulmate-dating/auth/internal/adapters/jwt"
	"github.com/soulmate-dating/auth/internal/adapters/postgres"
	"github.com/soulmate-dating/auth/internal/hash"
	"github.com/soulmate-dating/auth/internal/models"
)

type App interface {
	SignUp(ctx context.Context, credentials models.LoginCredentials) (*models.Token, error)
	Login(ctx context.Context, credentials models.LoginCredentials) (*models.Token, error)
	Refresh(ctx context.Context, token string) (*models.Token, error)
	Logout(ctx context.Context, token string) (string, error)
	Validate(ctx context.Context, token string) (string, error)
}

type Repository interface {
	CreateUser(ctx context.Context, p *models.User) (uuid.UUID, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type TransactionManager interface {
	RunInTx(ctx context.Context, f func(ctx context.Context) error) error
}

type Application struct {
	validate   *validator.Validate
	repository Repository
	jwtWrapper jwt.Wrapper
	txManager  TransactionManager
}

func (a *Application) Validate(_ context.Context, token string) (string, error) {
	err := a.validate.Var(token, "jwt")
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}
	claims, err := a.jwtWrapper.ValidateToken(token)
	if err != nil {
		return "", err
	}
	return claims.Id.String(), nil
}

func (a *Application) Logout(_ context.Context, token string) (string, error) {
	claims, err := a.jwtWrapper.ValidateToken(token)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}
	return claims.Id.String(), nil
}

func (a *Application) SignUp(ctx context.Context, credentials models.LoginCredentials) (token *models.Token, err error) {
	err = a.validate.Struct(credentials)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password: %w", err)
	}
	var user *models.User
	err = a.txManager.RunInTx(ctx, func(ctx context.Context) error {
		user, err = a.signup(ctx, credentials)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to signup: %w", err)
	}
	return a.generateTokenForUser(user)
}

func (a *Application) signup(ctx context.Context, credentials models.LoginCredentials) (*models.User, error) {
	user, err := a.repository.GetUserByEmail(ctx, credentials.Email)
	if err == nil {
		return nil, models.ErrAlreadyExists
	}

	user = &models.User{
		ID:       models.NewUUID(),
		Email:    credentials.Email,
		Password: hash.HashPassword(credentials.Password),
	}

	id, err := a.repository.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	user.ID = id
	return user, nil
}

func (a *Application) Login(ctx context.Context, credentials models.LoginCredentials) (*models.Token, error) {
	err := a.validate.Struct(credentials)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password: %w", err)
	}
	user, err := a.repository.GetUserByEmail(ctx, credentials.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	match := hash.CheckPasswordHash(credentials.Password, user.Password)
	if !match {
		return nil, models.ErrWrongPassword
	}

	return a.generateTokenForUser(user)
}

func (a *Application) Refresh(ctx context.Context, token string) (*models.Token, error) {
	err := a.validate.Var(token, "jwt")
	if err != nil {
		return nil, err
	}
	claims, err := a.jwtWrapper.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	user, err := a.repository.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return nil, err
	}

	return a.generateTokenForUser(user)
}

func (a *Application) generateTokenForUser(user *models.User) (*models.Token, error) {
	accessToken, err := a.jwtWrapper.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	refreshToken, err := a.jwtWrapper.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &models.Token{Id: user.ID, AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func NewApp(conn *pgxpool.Pool, jwt *jwt.Wrapper) App {
	pool := postgres.NewPool(conn)
	repo := postgres.NewRepo(pool)
	return &Application{repository: repo, jwtWrapper: *jwt, txManager: pool, validate: validator.New()}
}
