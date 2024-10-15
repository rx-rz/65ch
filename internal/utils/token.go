package utils

import (
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(payload map[string]string, secret string) string {
	claims := jwt.MapClaims{}
	for key, value := range payload {
		claims[key] = value
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}
