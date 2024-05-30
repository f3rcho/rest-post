package models

import "github.com/golang-jwt/jwt"

type AppClaims struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}
