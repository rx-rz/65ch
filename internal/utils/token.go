package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lucsky/cuid"
	"strings"
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

func DecodeToken[T jwt.Claims](tokenString string, secret string, claims T) (T, error) {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return claims, err
	}
	if _, ok := token.Claims.(T); ok && token.Valid {
		return claims, nil
	}
	return claims, fmt.Errorf("invalid token provided")
}

func GenerateResetToken() (string, time.Time) {
	parts := []string{cuid.New(), cuid.New()}
	token := strings.Join(parts, "")[:32]
	expiry := time.Now().Add(time.Minute * 15)
	return token, expiry
}
