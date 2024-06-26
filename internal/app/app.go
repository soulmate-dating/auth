package app

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/soulmate-dating/auth/internal/adapters/jwt"
	"github.com/soulmate-dating/auth/internal/adapters/postgres"
	"github.com/soulmate-dating/auth/internal/config"
	"github.com/soulmate-dating/auth/internal/domain"
	"github.com/soulmate-dating/auth/internal/hash"
	"log"
)

const jwtTag = "jwt"

type App interface {
	SignUp(ctx context.Context, credentials domain.LoginCredentials) (*domain.Token, error)
	Login(ctx context.Context, credentials domain.LoginCredentials) (*domain.Token, error)
	Refresh(ctx context.Context, token string) (*domain.Token, error)
	Logout(ctx context.Context, token string) (string, error)
	Validate(ctx context.Context, token string) (string, error)
}

type Repository interface {
	CreateUser(ctx context.Context, p *domain.User) (uuid.UUID, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
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
	err := a.validate.Var(token, jwtTag)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}
	claims, err := a.jwtWrapper.ValidateAccessToken(token)
	if err != nil {
		return "", err
	}
	return claims.Id.String(), nil
}

func (a *Application) Logout(_ context.Context, token string) (string, error) {
	claims, err := a.jwtWrapper.ValidateAccessToken(token)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}
	return claims.Id.String(), nil
}

func (a *Application) SignUp(ctx context.Context, credentials domain.LoginCredentials) (token *domain.Token, err error) {
	err = a.validate.Struct(credentials)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password: %w", err)
	}
	var user *domain.User
	err = a.txManager.RunInTx(ctx, func(ctx context.Context) error {
		user, err = a.signup(ctx, credentials)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to signup: %w", err)
	}
	return a.generateTokenForUser(user)
}

func (a *Application) signup(ctx context.Context, credentials domain.LoginCredentials) (*domain.User, error) {
	_, err := a.repository.GetUserByEmail(ctx, credentials.Email)
	if err == nil {
		return nil, domain.ErrAlreadyExists
	}

	user := &domain.User{
		ID:       domain.NewUUID(),
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

func (a *Application) Login(ctx context.Context, credentials domain.LoginCredentials) (*domain.Token, error) {
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
		return nil, domain.ErrWrongPassword
	}

	return a.generateTokenForUser(user)
}

func (a *Application) Refresh(ctx context.Context, token string) (*domain.Token, error) {
	err := a.validate.Var(token, jwtTag)
	if err != nil {
		return nil, err
	}
	claims, err := a.jwtWrapper.ValidateRefreshToken(token)
	if err != nil {
		return nil, err
	}
	user, err := a.repository.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return a.generateTokenForUser(user)
}

func (a *Application) generateTokenForUser(user *domain.User) (*domain.Token, error) {
	accessToken, err := a.jwtWrapper.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	refreshToken, err := a.jwtWrapper.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &domain.Token{Id: user.ID, AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func New(ctx context.Context, cfg config.Config) App {
	conn, err := postgres.Connect(ctx, postgres.Config{
		Host:              cfg.Postgres.Host,
		Port:              cfg.Postgres.Port,
		User:              cfg.Postgres.User,
		Password:          cfg.Postgres.Password,
		DBName:            cfg.Postgres.Database,
		SSLMode:           cfg.Postgres.SSLMode,
		ConnectionTimeout: cfg.Postgres.ConnectionTimeout,
	})
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	wrapper := jwt.NewWrapper(
		cfg.JWT.SecretKey,
		cfg.JWT.RefreshSecretKey,
		cfg.JWT.Issuer,
		cfg.JWT.AccessExpirationHours,
		cfg.JWT.RefreshExpirationHours,
	)
	pool := postgres.NewPool(conn)
	repo := postgres.NewRepo(pool)
	return &Application{repository: repo, jwtWrapper: *wrapper, txManager: pool, validate: validator.New()}
}
