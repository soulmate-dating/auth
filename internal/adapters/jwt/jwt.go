package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/soulmate-dating/auth/internal/domain"
)

type Wrapper struct {
	SecretKey              string
	RefreshSecretKey       string
	Issuer                 string
	AccessTokenExpiration  time.Duration
	RefreshTokenExpiration time.Duration
}

func NewWrapper(
	issuer, secretKey, refreshSecretKey string,
	accessExpirationHours, refreshExpirationHours time.Duration,
) *Wrapper {
	return &Wrapper{
		SecretKey:              secretKey,
		RefreshSecretKey:       refreshSecretKey,
		Issuer:                 issuer,
		AccessTokenExpiration:  accessExpirationHours,
		RefreshTokenExpiration: refreshExpirationHours,
	}
}

func (w *Wrapper) GenerateAccessToken(user *domain.User) (string, error) {
	return w.generateToken(user, w.AccessTokenExpiration, w.SecretKey)
}

func (w *Wrapper) GenerateRefreshToken(user *domain.User) (string, error) {
	return w.generateToken(user, w.RefreshTokenExpiration, w.RefreshSecretKey)
}

func (w *Wrapper) generateToken(user *domain.User, expiration time.Duration, secretKey string) (signedToken string, err error) {
	claims := &domain.Claims{
		Id:    user.ID,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(expiration).Unix(),
			Issuer:    w.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (w *Wrapper) ValidateAccessToken(signedToken string) (claims *domain.Claims, err error) {
	return w.validateToken(signedToken, w.SecretKey)
}

func (w *Wrapper) ValidateRefreshToken(signedToken string) (claims *domain.Claims, err error) {
	return w.validateToken(signedToken, w.RefreshSecretKey)
}

func (w *Wrapper) validateToken(signedToken, secretKey string) (claims *domain.Claims, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&domain.Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*domain.Claims)
	if !ok {
		return nil, domain.ErrFailedToParseClaims
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, domain.ErrExpiredToken
	}

	return claims, nil
}
