package model

import "github.com/golang-jwt/jwt/v5"

var JwtKey = []byte("secret-key")

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}