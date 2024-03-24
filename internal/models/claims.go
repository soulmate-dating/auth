package models

import "github.com/golang-jwt/jwt"

type Claims struct {
	jwt.StandardClaims
	Id    string
	Email string
}
