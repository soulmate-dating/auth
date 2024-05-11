package jwt

import (
	"fmt"
	"github.com/soulmate-dating/auth/internal/models"
	"time"

	"github.com/golang-jwt/jwt"
)

type Wrapper struct {
	SecretKey              string
	Issuer                 string
	AccessExpirationHours  int64
	RefreshExpirationHours int64
}

func NewWrapper(
	secretKey, issuer string,
	accessExpirationHours, refreshExpirationHours int64,
) *Wrapper {
	return &Wrapper{
		SecretKey:              secretKey,
		Issuer:                 issuer,
		AccessExpirationHours:  accessExpirationHours,
		RefreshExpirationHours: refreshExpirationHours,
	}
}

func (w *Wrapper) GenerateAccessToken(user *models.User) (string, error) {
	return w.generateToken(user, w.AccessExpirationHours)
}

func (w *Wrapper) GenerateRefreshToken(user *models.User) (string, error) {
	return w.generateToken(user, w.RefreshExpirationHours)
}

func (w *Wrapper) generateToken(user *models.User, expirationHours int64) (signedToken string, err error) {
	claims := &models.Claims{
		Id:    user.ID,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(expirationHours)).Unix(),
			Issuer:    w.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(w.SecretKey))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (w *Wrapper) ValidateToken(signedToken string) (claims *models.Claims, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&models.Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(w.SecretKey), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", models.ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok {
		return nil, models.ErrFailedToParseClaims
	}

	return claims, nil
}
