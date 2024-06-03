package utils

import (
	"net/http"
	"strings"

	"github.com/f3rcho/rest-posts/models"
	"github.com/f3rcho/rest-posts/server"
	"github.com/golang-jwt/jwt"
)

func GetClaims(s server.Server, w http.ResponseWriter, r *http.Request) (*models.AppClaims, error) {
	tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
	token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.Config().JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*models.AppClaims)

	if ok && token.Valid {
		return claims, err
	}

	return nil, err
}
