package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateToken(payload map[string]string, secret string) (string, error) {
	expirationTime := time.Now().Add(time.Minute * 15).Unix()
	claims := jwt.MapClaims{
		"exp": expirationTime,
	}
	for key, value := range payload {
		claims[key] = value
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
