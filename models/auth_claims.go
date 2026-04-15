package models

import "github.com/golang-jwt/jwt/v5"

type AuthClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}
