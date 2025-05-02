package server

import (
	"github.com/golang-jwt/jwt/v5"
)

var JwtKey = []byte("secret_key")

type Claims struct {
	UserID int    `json:"user_id"`
	Login  string `json:"login"`
	jwt.RegisteredClaims
}
